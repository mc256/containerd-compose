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

package app

import (
	"containerd-compose/version"

	cmdUp "containerd-compose/cmd/up"
	cmdVersion "containerd-compose/cmd/version"

	"fmt"
	"github.com/urfave/cli/v2"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, c.App.Version)
	}
}

func New() *cli.App {
	app := cli.NewApp()

	app.Name = "containerd-compose"
	app.Version = version.Version
	app.Description = `
containerd-compose is a tool to run multi-container applications with containerd.
It searches containerd-compose.yml and docker-compose.yml in current working directory by default.
`
	app.Usage = `
                    __        _                     __                                               
  _________  ____  / /_____ _(_)___  ___  _________/ /     _________  ____ ___  ____  ____  ________ 
 / ___/ __ \/ __ \/ __/ __  / / __ \/ _ \/ ___/ __  /_____/ ___/ __ \/ __  __ \/ __ \/ __ \/ ___/ _ \
/ /__/ /_/ / / / / /_/ /_/ / / / / /  __/ /  / /_/ /_____/ /__/ /_/ / / / / / / /_/ / /_/ (__  )  __/
\___/\____/_/ /_/\__/\__,_/_/_/ /_/\___/_/   \__,_/      \___/\____/_/ /_/ /_/ .___/\____/____/\___/ 
                                                                            /_/
containerd-compose
`

	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "compose file",
			DefaultText: "(containerd|docker)-compose.y(a)ml",
			HasBeenSet:  true,
		},
		&cli.StringFlag{
			Name:        "host",
			Aliases:     []string{"H", "s"},
			Usage:       "Containerd daemon socket to connect to",
			DefaultText: "/run",
			HasBeenSet:  true,
		},
		&cli.StringFlag{
			Name:        "namespace",
			Aliases:     []string{"ns"},
			Usage:       "Containerd namespaces",
			DefaultText: "default",
			HasBeenSet:  true,
		},
	}
	app.Commands = append([]*cli.Command{
		cmdVersion.Command(),
		cmdUp.Command(),
	})

	return app
}
