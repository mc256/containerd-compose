package up

import (
	"github.com/urfave/cli/v2"
)

func Action(context *cli.Context) error {
	//TODO: bring up the container

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
