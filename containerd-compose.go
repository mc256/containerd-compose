package main

import (
	cmdApp "containerd-compose/cmd/app"
	"fmt"
	"os"
)

func main() {
	app := cmdApp.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "containerd-compose: %v\n", err)
		os.Exit(1)
	}
}
