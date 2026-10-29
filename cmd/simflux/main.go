package main

import (
	"os"

	"github.com/Faithful001/simflux.git/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}