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
	"containerd-compose/operations/composer"
	"containerd-compose/version"

	cmdDown "containerd-compose/cmd/down"
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
			Usage:       "compose file name",
			DefaultText: "(containerd|docker)-compose.y(a)ml",
		},
		&cli.StringFlag{
			Name:        "project-name",
			Aliases:     []string{"p"},
			Usage:       "Specify an alternate project name",
			DefaultText: "directory name",
		},
		&cli.StringFlag{
			Name:        "host",
			Aliases:     []string{"H", "s"},
			Usage:       "Containerd daemon socket to connect to",
			DefaultText: "/run/containerd/containerd.sock",
		},
		&cli.StringFlag{
			Name:        "env-file",
			Aliases:     []string{"env", "e"},
			Usage:       "Specify an alternate environment file",
			DefaultText: "./.env",
		},
		&cli.StringFlag{
			Name:        "namespace",
			Aliases:     []string{"ns"},
			Usage:       "Containerd namespaces",
			DefaultText: "default",
		},
		&cli.StringFlag{
			Name:        "volume",
			Aliases:     []string{"vol"},
			Usage:       "Default path to storage volumes",
			DefaultText: "./volumes",
		},
	}
	app.Commands = append([]*cli.Command{
		cmdVersion.Command(),
		cmdUp.Command(ParseContext),
		cmdDown.Command(ParseContext),
	})

	return app
}

func ParseContext(context *cli.Context) ([]composer.Option, error) {
	var opts []composer.Option

	if context.Bool("debug") {
		opts = append(opts, composer.WithDebugMode())
	}

	if composeFile := context.String("file"); composeFile != "" {
		opts = append(opts, composer.WithComposeFile(composeFile))
	}

	if projectName := context.String("project-name"); projectName != "" {
		opts = append(opts, composer.WithProjectName(projectName))
	}

	if host := context.String("host"); host != "" {
		opts = append(opts, composer.WithContainerdSocketFile(host))
	}

	if envFile := context.String("env-file"); envFile != "" {
		opts = append(opts, composer.WithEnvFile(envFile))
	}

	if namespace := context.String("namespace"); namespace != "" {
		opts = append(opts, composer.WithNamespace(namespace))
	}

	if volume := context.String("volume"); volume != "" {
		opts = append(opts, composer.WithVolumeBase(volume))
	}

	return opts, nil
}
