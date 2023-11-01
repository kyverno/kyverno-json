package scan

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
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
	cmd.Flags().StringSliceVar(&command.preprocessors, "pre-process", nil, "JMESPath expression used to pre process payload")
	cmd.Flags().StringSliceVar(&command.policies, "policy", nil, "Path to kyverno-json policies")
	cmd.Flags().StringSliceVar(&command.selectors, "labels", nil, "Labels selectors for policies")
	cmd.Flags().StringVar(&command.identifier, "identifier", "", "JMESPath expression used to identify a resource")
	return cmd
}
