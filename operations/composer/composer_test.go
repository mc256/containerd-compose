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
	"fmt"
	"testing"
)

func Test_getImageFullName(t *testing.T) {
	out, err := getFullImageName("nextcloud", "docker.io/library")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(out)
}

func Test_ComposerLoaderProvideComposeFile(t *testing.T) {
	var opts []Option

	//opts = append(opts, WithComposeFile("./test.yml"))
	opts = append(opts, WithComposeFile("./docker-compose.yml"))
	var compose *ComposeFile
	var err error
	if compose, err = LoadFile(opts...); err != nil {
		t.Error(err)
	}
	if compose.Version != "2" {
		t.Error("version error")
	}
	fmt.Println(compose)
}

func Test_ComposerLoaderNonProvideComposeFile(t *testing.T) {
	var opts []Option

	opts = append(opts, WithComposeFile("./xxx-compose.yml"), WithEnvFile("env.txt"))
	var compose *ComposeFile
	var err error
	if compose, err = LoadFile(opts...); err != nil {
		fmt.Println(err)
	}
	if compose == nil {
		fmt.Println("OK")
	}
}

func Test_LoadContainerd(t *testing.T) {
	var opts []Option
	var compose *ComposeFile
	var err error
	if compose, err = LoadFile(opts...); err != nil {
		t.Error(err)
	}
	if err := LaunchApplication(compose, opts...); err != nil {
		t.Error(err)
	}
}
