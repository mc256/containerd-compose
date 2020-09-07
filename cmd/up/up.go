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

package up

import (
	"containerd-compose/operations/composer"
	"fmt"
	"github.com/urfave/cli/v2"
)

func Action(context *cli.Context) error {
	var opts []composer.Option

	// parse containerd-compose file
	var compose *composer.ComposeFile
	var err error
	if compose, err = composer.LoadFile(opts...); err != nil {
		return err
	}

	//TODO: bring up the container
	fmt.Println(*compose)

	return nil
}

func Command() *cli.Command {
	cmd := cli.Command{
		Name:   "up",
		Usage:  "Builds, (re)creates, starts, and attached to containers for a service.",
		Action: Action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "detach",
				Aliases: []string{"d"},
				Usage:   "run application in background.",
			},
		},
	}
	return &cmd
}
