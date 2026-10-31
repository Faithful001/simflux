package main

import (
	"os"

	"github.com/Faithful001/simflux.git/internal/cli"
	"github.com/Faithful001/simflux.git/internal/config/db"
)

func main() {
	db.InitDB()

	if err := db.AutoMigrate(); err != nil {
		os.Exit(1)
	}


	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}