package version

import (
	"fmt"

	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/kyverno/kyverno-json/pkg/version"
	"github.com/spf13/cobra"
)

func Command(parent *cobra.Command) *cobra.Command {
	doc := command.New(
		parent,
		command.WithDescription("Prints the version informations."),
		command.WithExample("Print version infos", "version"),
	)
	return &cobra.Command{
		Use:          "version",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", version.Version())
			fmt.Fprintf(cmd.OutOrStdout(), "Time: %s\n", version.Time())
			fmt.Fprintf(cmd.OutOrStdout(), "Git commit ID: %s\n", version.Hash())
			return nil
		},
	}
}
