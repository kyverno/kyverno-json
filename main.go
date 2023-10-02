package main

import (
	"os"

	"github.com/kyverno/kyverno-json/pkg/commands"
)

func main() {
	root := commands.NewRootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
