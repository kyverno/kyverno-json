package jp

import (
	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/kyverno/kyverno-json/pkg/commands/jp/function"
	"github.com/kyverno/kyverno-json/pkg/commands/jp/parse"
	"github.com/kyverno/kyverno-json/pkg/commands/jp/query"
	"github.com/spf13/cobra"
)

func Command(parents ...string) *cobra.Command {
	doc := command.New(
		command.WithParents(parents...),
		command.WithDescription("Provides a command-line interface to JMESPath, enhanced with custom functions."),
		command.WithExample("List functions", "jp function"),
		command.WithExample("Evaluate query", "jp query -i object.yaml 'request.object.metadata.name | truncate(@, `9`)'"),
		command.WithExample("Parse expression", "jp parse 'request.object.metadata.name | truncate(@, `9`)'"),
	)
	cmd := &cobra.Command{
		Use:          "jp",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		function.Command(append(parents, "jp")...),
		parse.Command(append(parents, "jp")...),
		query.Command(append(parents, "jp")...),
	)
	return cmd
}
