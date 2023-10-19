package commands

import (
	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/kyverno/kyverno-json/pkg/commands/docs"
	"github.com/kyverno/kyverno-json/pkg/commands/jp"
	"github.com/kyverno/kyverno-json/pkg/commands/playground"
	"github.com/kyverno/kyverno-json/pkg/commands/scan"
	"github.com/kyverno/kyverno-json/pkg/commands/serve"
	"github.com/kyverno/kyverno-json/pkg/commands/version"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	doc := command.New(
		command.WithDescription(
			"kyverno-json is a CLI tool to apply policies to json resources.",
		),
	)
	cmd := &cobra.Command{
		Use:          "kyverno-json",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		docs.Command("kyverno-json"),
		jp.Command("kyverno-json"),
		playground.Command(),
		scan.Command(),
		serve.Command("kyverno-json"),
		version.Command("kyverno-json"),
	)
	return cmd
}
