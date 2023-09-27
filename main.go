package main

import (
	"os"

	"github.com/eddycharly/tf-kyverno/pkg/commands"
)

func main() {
	root := commands.NewRootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
