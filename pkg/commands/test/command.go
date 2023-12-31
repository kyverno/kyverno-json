package test

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var command options
	cmd := &cobra.Command{
		Use:          "test",
		Short:        "test",
		Long:         "Run tests from a local filesystem",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE:         command.run,
	}

	cmd.Flags().StringVarP(&command.fileName, "file-name", "f", "kyverno-test.yaml", "Test filename")
	cmd.Flags().BoolVar(&command.removeColor, "remove-color", false, "Remove any color from output")
	return cmd
}
