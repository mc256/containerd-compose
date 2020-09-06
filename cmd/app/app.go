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
			DefaultText: "containerd-compose.yml, docker-compose.yml",
			HasBeenSet:  true,
		},
	}
	app.Commands = append([]*cli.Command{
		cmdVersion.Command(),
		cmdUp.Command(),
	})

	return app
}
