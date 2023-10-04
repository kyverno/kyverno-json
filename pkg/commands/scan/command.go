package scan

import (
	"github.com/spf13/cobra"
)

func Command(root *cobra.Command) *cobra.Command {
	var command options
	cmd := &cobra.Command{
		Use:          "scan",
		Short:        "scan",
		Long:         "Apply policies to json resources",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         command.run,
	}
	cmd.Flags().StringVar(&command.payload, "payload", "", "Path to payload (json or yaml file)")
	cmd.Flags().StringSliceVar(&command.preprocessors, "pre-process", nil, "JmesPath expression used to pre process payload")
	cmd.Flags().StringSliceVar(&command.policies, "policy", nil, "Path to kyverno-json policies")
	return cmd
}
