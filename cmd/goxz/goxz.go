package main

import (
	"os"

	cmd "github.com/aberestyak/goxz/pkg/cmd"
)

func main() {
	command := cmd.NewCmdGoXz()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
