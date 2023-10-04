package jp

import (
	"github.com/kyverno/kyverno-json/pkg/commands/jp/function"
	"github.com/kyverno/kyverno-json/pkg/commands/jp/parse"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "jp",
		Short:        "Provides a command-line interface to JMESPath, enhanced with custom functions.",
		Long:         "Provides a command-line interface to JMESPath, enhanced with custom functions.",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		function.Command(),
		parse.Command(),
		// query.Command(),
	)
	return cmd
}
