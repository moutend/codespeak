package main

import (
	"os"

	"github.com/moutend/codespeak/internal/cli"
)

func main() {
	cli.RootCommand.SetOutput(os.Stdout)

	if err := cli.RootCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}
