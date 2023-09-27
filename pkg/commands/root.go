package commands

import (
	"fmt"

	"github.com/eddycharly/tf-kyverno/pkg/policy"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/output/pluralize"
	"github.com/spf13/cobra"
)

type command struct {
	plan     string
	policies []string
}

func (c *command) Run(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Loading policies ...")
	policies, err := policy.Load(c.policies...)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "-", len(policies), pluralize.Pluralize(len(policies), "policy", "policies"), "loaded")
	fmt.Fprintln(out, "Loading plan ...")
	// TODO
	fmt.Fprintln(out, "Running ...")
	// TODO
	fmt.Fprintln(out, "Done")
	// TODO
	return nil
}

func NewRootCommand() *cobra.Command {
	var command command
	cmd := &cobra.Command{
		Use:          "tf-kyverno",
		Short:        "tf-kyverno",
		Long:         "tf-kyverno is a CLI tool to apply policies to terraform plans",
		Args:         cobra.NoArgs,
		RunE:         command.Run,
		SilenceUsage: true,
	}
	cmd.Flags().StringVar(&command.plan, "plan", "", "Path to terraform plan file (in json format)")
	cmd.Flags().StringSliceVar(&command.policies, "policy", nil, "Path to tf-kyverno policies")
	return cmd
}
