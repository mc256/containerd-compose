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
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

type ContextParser func(context *cli.Context) ([]Option, error)
type Option func(*options)

type options struct {
	isDebug          bool
	isDetach         bool
	isForced         bool
	composeFile      []string
	projectName      string
	projectDir       string
	envFile          string
	containerdSocket string
	namespace        string
	defaultRegistry  string
}

func WithDebugMode() Option {
	return func(opts *options) {
		opts.isDebug = true
	}
}

func WithDetachMode() Option {
	return func(opts *options) {
		opts.isDetach = true
	}
}

func WithForce() Option {
	return func(opts *options) {
		opts.isForced = true
	}
}

func WithComposeFile(composeFile string) Option {
	return func(opts *options) {
		opts.composeFile = []string{
			composeFile,
		}
	}
}

func WithProjectName(projectName string) Option {
	return func(opts *options) {
		opts.projectName = projectName
	}
}

func WithProjectDir(projectDir string) Option {
	return func(opts *options) {
		opts.projectDir = projectDir
	}
}

func WithEnvFile(envFile string) Option {
	return func(opts *options) {
		opts.envFile = envFile
	}
}

func WithContainerdSocketFile(socket string) Option {
	return func(opts *options) {
		opts.containerdSocket = socket
	}
}

func WithNamespace(namespace string) Option {
	return func(opts *options) {
		opts.namespace = namespace
	}
}

func WithDefaultImageRegistry(registry string) Option {
	return func(opts *options) {
		opts.defaultRegistry = registry
	}
}

func ContextToOptions(context *cli.Context, contextParsers ...ContextParser) ([]Option, error) {
	var opts []Option
	for _, p := range contextParsers {
		temp, err := p(context)
		if err != nil {
			return nil, err
		}
		opts = append(opts, temp...)
	}
	return opts, nil
}

func parseOptions(opts *[]Option) *options {
	opt := options{}
	for _, o := range *opts {
		o(&opt)
	}
	if opt.envFile == "" {
		opt.envFile = ".env"
	}
	if opt.projectDir == "" {
		opt.projectDir = "./"
	}
	if len(opt.composeFile) == 0 {
		opt.composeFile = []string{
			"./containerd-compose.yml",
			"./docker-compose.yml",
			"./containerd-compose.yaml",
			"./docker-compose.yaml",
		}
	}

	if opt.projectName == "" {
		dir, err := os.Getwd()
		if err != nil {
			opt.projectName = "void"
		} else {
			split := strings.Split(dir, "/")
			opt.projectName = split[len(split)-1]
		}
	}

	if opt.containerdSocket == "" {
		opt.containerdSocket = "/run/containerd/containerd.sock"
	}

	if opt.namespace == "" {
		opt.namespace = "default"
	}

	if opt.defaultRegistry == "" {
		opt.defaultRegistry = "docker.io/library"
	}

	return &opt
}
