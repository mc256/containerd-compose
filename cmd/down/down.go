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

package down

import (
	"containerd-compose/operations/composer"
	"github.com/urfave/cli/v2"
)

func Action(context *cli.Context, contextParsers ...composer.ContextParser) error {
	opts, err := composer.ContextToOptions(context, contextParsers...)
	if err != nil {
		return err
	}

	var compose *composer.ComposeFile
	if compose, err = composer.LoadFile(opts...); err != nil {
		return err
	}

	if err := composer.StopApplication(compose, opts...); err != nil {
		return err
	}

	return nil
}

func Command(parsers ...composer.ContextParser) *cli.Command {
	cmd := cli.Command{
		Name:  "up",
		Usage: "Builds, (re)creates, starts, and attached to containers for a service.",
		Action: func(c *cli.Context) error {
			return Action(c, append(parsers, ParseContext)...)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "force to exit",
			},
		},
	}
	return &cmd
}

func ParseContext(context *cli.Context) ([]composer.Option, error) {
	var opts []composer.Option

	if context.Bool("force") {
		opts = append(opts, composer.WithForce())
	}

	return opts, nil
}
