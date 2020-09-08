/*
   Copyright The containerd-compose Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

/*
   file created by mc256.com in 2020
*/

package composer

import (
	"containerd-compose/logger"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/snapshots"
	"github.com/joho/godotenv"
	"github.com/opencontainers/runtime-spec/specs-go"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func init() {

}

func getFullImageName(s string, base string) (string, error) {
	split := strings.Split(s, "/")
	size := len(split)
	if size == 0 {
		return "", errors.New("not a valid image name")
	} else {
		tagSplit := strings.Split(split[size-1], ":")
		if len(tagSplit) == 1 {
			s = s + ":latest"
		}
		if size == 1 {
			return path.Join(base, s), nil
		} else {
			return s, nil
		}
	}
}

func LoadFile(opts ...Option) (*ComposeFile, error) {
	// Options
	var opt *options
	opt = parseOptions(&opts)

	// Read yaml file from local
	var buffer []byte
	var err error

	// try default value
	for _, d := range opt.composeFile {
		buffer, err = ioutil.ReadFile(d)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, errors.New("compose file not found")
	}

	// Get Environment Variables
	_ = godotenv.Load(opt.envFile)

	buffer = []byte(os.ExpandEnv(string(buffer)))

	// Parse Yaml file
	t := ComposeFile{}
	if err := yaml.Unmarshal(buffer, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func calculateVolumeHash(projectName string, serviceName string, dst string) (result string) {
	h := sha1.New()
	_, _ = io.WriteString(h, projectName)
	_, _ = io.WriteString(h, serviceName)
	_, _ = io.WriteString(h, dst)
	result = fmt.Sprintf("%x", h.Sum(nil))[:10]
	return
}

func createMountingDir(dir string) error {
	dirSplit := strings.Split(dir, "/")
	lastSplit := strings.Split(dirSplit[len(dirSplit)-1], ".")
	if len(lastSplit) > 1 {
		return nil
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func parseVolumesConfiguration(volume string, projectName string, serviceName string, opt *options) (mount *specs.Mount, err error) {
	sl := strings.Split(volume, ":")
	var src, dst string
	var mountOpts []string
	if len(sl) == 1 {
		dst = sl[0]
		src = filepath.Join(opt.volumeBase, calculateVolumeHash(projectName, serviceName, dst))
		_ = createMountingDir(src)
		mountOpts = []string{"rw"}
	} else {
		if len(sl) == 2 {
			mountOpts = []string{"rw"}
		} else {
			mountOpts = sl[2:]
		}

		src = sl[0]
		dst = sl[1]
		if strings.HasPrefix(src, "~/") {
			usr, _ := user.Current()
			src = filepath.Join(usr.HomeDir, src[2:])
		} else if strings.HasPrefix(src, "./") {
			pwd, _ := os.Getwd()
			src = filepath.Join(pwd, src[2:])
		} else if strings.HasPrefix(src, "/") {
			// happy
		} else {
			src = filepath.Join(opt.volumeBase, calculateVolumeHash(projectName, serviceName, dst))
		}
	}
	if err := createMountingDir(src); err != nil {
		return nil, err
	}
	//fmt.Println(mountOpts)
	mount = &specs.Mount{
		Destination: dst,
		Source:      src,
		Type:        "bind",
		Options:     append([]string{"rbind"}, mountOpts...),
		//Options:     []string{"rw"},
	}
	return mount, nil
}

func LaunchApplication(compose *ComposeFile, opts ...Option) error {
	// Options
	var opt *options
	opt = parseOptions(&opts)
	logger.N("project: %s --- <up>", opt.projectName)

	// Connect to containerd service
	client, err := containerd.New(opt.containerdSocket)
	if err != nil {
		return err
	}
	defer client.Close()
	logger.N("connected to containerd: %s", opt.containerdSocket)

	// Pull Image to Namespace
	ctx := namespaces.WithNamespace(context.Background(), opt.namespace)
	logger.N("using namespace: %s", opt.namespace)

	for serviceName, serviceConfig := range compose.Services {
		logger.N("building service: %s", serviceName)
		imageName, err := getFullImageName(serviceConfig.Image, opt.defaultRegistry)
		if err != nil {
			return err
		}

		logger.N("pulling image: %s", imageName)
		image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
		if err != nil {
			return err
		}

		// ---------------------------------------------------------------
		containerId := fmt.Sprintf("%s-%s", opt.projectName, serviceName)

		// Prepare Mounting
		logger.N("preparing mounting points: %s", containerId)
		var mounts []specs.Mount
		for _, v := range serviceConfig.Volumes {
			if mount, err := parseVolumesConfiguration(v, opt.projectName, serviceName, opt); err != nil {
				return err
			} else {
				mounts = append(mounts, *mount)
			}
		}

		for _, otherService := range serviceConfig.VolumesFrom {
			for _, v := range compose.Services[otherService].Volumes {
				if mount, err := parseVolumesConfiguration(v, opt.projectName, otherService, opt); err != nil {
					return err
				} else {
					mounts = append(mounts, *mount)
				}
			}
		}
		logger.N("mounting points: %v", mounts)

		// Create Container
		logger.N("creating container: %s", containerId)
		var spec []oci.SpecOpts
		spec = append(spec, oci.WithImageConfig(image))
		spec = append(spec, oci.WithEnv(serviceConfig.Environment))
		spec = append(spec, oci.WithMounts(mounts))
		spec = append(spec, oci.WithHostNamespace(specs.NetworkNamespace))
		spec = append(spec, oci.WithHostHostsFile)
		spec = append(spec, oci.WithHostResolvconf)

		//spec = append(spec, oci.WithPrivileged)
		//spec = append(spec, oci.WithAllDevicesAllowed)

		var snapshotOps []snapshots.Opt

		container, err := client.NewContainer(
			ctx,
			containerId,
			containerd.WithImage(image),
			containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", containerId), image, snapshotOps...),
			containerd.WithNewSpec(spec...),
		)

		if err != nil {
			return err
		}

		// Create Task
		var task containerd.Task
		logger.N("creating task: %s", containerId)
		if opt.isDetach {
			task, err = container.NewTask(ctx, cio.NewCreator())
			if err != nil {
				return err
			}
		} else {
			stdout := logger.NewStreamLogger(containerId, "out")
			stderr := logger.NewStreamLogger(containerId, "err")
			task, err = container.NewTask(ctx, cio.NewCreator(cio.WithStreams(nil, stdout, stderr)))
			if err != nil {
				return err
			}
		}

		// make sure we wait before calling start
		_, err = task.Wait(ctx)
		if err != nil {
			return err
		}

		// call start on the task to execute the redis server
		logger.N("launched: %s", containerId)
		if err := task.Start(ctx); err != nil {
			return err
		}
	}

	if !opt.isDetach {
		signalChan := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-signalChan
			fmt.Println()
			logger.N("received signal: %s", sig)
			done <- true
		}()
		<-done
		if err := StopApplication(compose, opts...); err != nil {
			return err
		}
	}

	return nil
}

func StopApplication(compose *ComposeFile, opts ...Option) error {

	// Options
	var opt *options
	opt = parseOptions(&opts)
	logger.N("project: %s --- <down>", opt.projectName)

	// Connect to containerd service
	client, err := containerd.New(opt.containerdSocket)
	if err != nil {
		return err
	}
	defer client.Close()
	logger.N("connected to containerd: %s", opt.containerdSocket)

	// Pull Image to Namespace
	ctx := namespaces.WithNamespace(context.Background(), opt.namespace)
	logger.N("using namespace: %s", opt.namespace)

	for k, _ := range compose.Services {
		containerId := fmt.Sprintf("%s-%s", opt.projectName, k)
		logger.N("stopping service: %s", containerId)
		// ---------------------------------------------------------------
		container, err := client.LoadContainer(ctx, containerId)
		if err != nil {
			logger.N("load container error: %s - %v", containerId, err)
		} else {
			task, err := container.Task(ctx, nil)
			if err != nil {
				logger.N("load task error: %s - %v", containerId, err)
			} else {
				exitStatus, err := task.Wait(ctx)
				if err != nil {
					logger.N("task wait error: %s - %v", containerId, err)
				}

				keyboardSignal := make(chan os.Signal, 1)
				done := make(chan bool, 1)
				signal.Notify(keyboardSignal, syscall.SIGINT, syscall.SIGTERM)

				var task_killed = false
				go func() {
					time.Sleep(10 * time.Second)
					if !task_killed {
						logger.N("killing task %s takes longer than 10 seconds, press Ctrl+C to force exit", containerId)
					}
				}()

				if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
					logger.N("kill task error: %s - %v", containerId, err)
					done <- true
				} else {
					done <- true
				}

				go func() {
					_ = <-keyboardSignal
					if !task_killed {
						_ = task.Kill(ctx, syscall.SIGKILL)
						fmt.Println()
						logger.N("sent SIGKILL: %s", containerId)
						done <- true
					}
				}()

				<-done
				status := <-exitStatus

				code, _, err := status.Result()
				logger.N("task exited (%d) with error: %v", code, err)
				task_killed = true

				if _, err := task.Delete(ctx); err != nil {
					logger.N("delete task error: %s - %v", containerId, err)
				}
			}
			if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
				logger.N("delete container error: %s - %v", containerId, err)
			}
		}
	}
	logger.N("project '%s' down", opt.projectName)
	return nil
}
