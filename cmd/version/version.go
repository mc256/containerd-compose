package version

import (
	"containerd-compose/version"
	"fmt"
	"github.com/urfave/cli/v2"
)

func Action(context *cli.Context) error {
	fmt.Printf("containerd-compose version %s", version.Version)
	return nil
}

func Command() *cli.Command {
	cmd := cli.Command{
		Name:   "version",
		Usage:  "show version information",
		Action: Action,
	}
	return &cmd
}
