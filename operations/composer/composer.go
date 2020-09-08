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
	"github.com/joho/godotenv"
	"github.com/opencontainers/runtime-spec/specs-go"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
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
	_ = godotenv.Load(".env")

	buffer = []byte(os.ExpandEnv(string(buffer)))

	// Parse Yaml file
	t := ComposeFile{}
	if err := yaml.Unmarshal(buffer, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func calculateVolumeHash(containerId string, dst string) (result string) {
	h := sha1.New()
	_, _ = io.WriteString(h, containerId)
	_, _ = io.WriteString(h, dst)
	result = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func createMountingDir(dir string) error {
	dirSplit := strings.Split(dir, "/")
	lastSplit := strings.Split(dirSplit[len(dirSplit)-1], ".")
	if len(lastSplit) >= 1 {
		return nil
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
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

	for k, s := range compose.Services {
		logger.N("building service: %s", k)
		imageName, err := getFullImageName(s.Image, opt.defaultRegistry)
		if err != nil {
			return err
		}

		logger.N("pulling image: %s", imageName)
		image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
		if err != nil {
			return err
		}

		// ---------------------------------------------------------------
		containerId := fmt.Sprintf("%s-%s", opt.projectName, k)

		// Prepare Mounting
		logger.N("preparing mounting points: %s", containerId)
		var mounts []specs.Mount
		for _, v := range s.Volumes {
			sl := strings.Split(v, ":")
			var src, dst string
			var mountOpts []string
			if len(sl) == 1 {
				dst = sl[0]
				src = filepath.Join(opt.volumeBase, calculateVolumeHash(containerId, dst+k))
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
					src = filepath.Join(opt.volumeBase, calculateVolumeHash(containerId, dst))
				}
			}
			if err := createMountingDir(src); err != nil {
				return err
			}
			mounts = append(mounts, specs.Mount{
				Destination: dst,
				Source:      src,
				Type:        "bind",
				Options:     append([]string{"rbind"}, mountOpts...),
			})
		}

		// TODO: Volumes Form

		// Create Container
		logger.N("creating container: %s", containerId)
		var spec []oci.SpecOpts
		spec = append(spec, oci.WithImageConfig(image))
		spec = append(spec, oci.WithEnv(s.Environment))
		spec = append(spec, oci.WithMounts(mounts))
		fmt.Println(mounts)

		container, err := client.NewContainer(
			ctx,
			containerId,
			containerd.WithImage(image),
			containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", containerId), image),
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
				if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
					logger.N("kill task error: %s - %v", containerId, err)
				}
				status := <-exitStatus
				code, _, err := status.Result()
				logger.N("task exited (%d) with error: %v", code, err)

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
