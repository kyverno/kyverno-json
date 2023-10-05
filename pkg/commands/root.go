package commands

import (
	"github.com/kyverno/kyverno-json/pkg/commands/docs"
	"github.com/kyverno/kyverno-json/pkg/commands/jp"
	"github.com/kyverno/kyverno-json/pkg/commands/scan"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kyverno-json",
		Short:        "kyverno-json",
		Long:         "kyverno-json is a CLI tool to apply policies to json resources",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		docs.Command(cmd),
		jp.Command(cmd),
		scan.Command(),
	)
	return cmd
}
