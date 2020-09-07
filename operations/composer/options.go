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

type Option func(*options)

type options struct {
	composeFile string
	projectDir  string
	envFile     string
}

func WithComposeFile(composeFile string) Option {
	return func(opts *options) {
		opts.composeFile = composeFile
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
