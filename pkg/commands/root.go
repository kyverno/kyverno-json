package commands

import (
	"errors"
	"fmt"

	"github.com/eddycharly/tf-kyverno/pkg/engine"
	"github.com/eddycharly/tf-kyverno/pkg/plan"
	"github.com/eddycharly/tf-kyverno/pkg/policy"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/output/pluralize"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	plan, err := plan.Load(c.plan)
	if err != nil {
		return err
	}
	resources, ok, err := unstructured.NestedSlice(plan, "planned_values", "root_module", "resources")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to find resources in the plan")
	}
	fmt.Fprintln(out, "-", len(resources), pluralize.Pluralize(len(resources), "resource", "resources"), "loaded")
	fmt.Fprintln(out, "Running ...")
	// TODO
	for _, resource := range resources {
		resourceName, _, _ := unstructured.NestedString(resource.(map[string]interface{}), "address")
		for _, policy := range policies {
			for _, rule := range policy.Spec.Rules {
				match, exclude := engine.MatchExclude(rule.MatchResources, rule.ExcludeResources, resource)
				if match && !exclude {
					fmt.Fprintln(out, "-", policy.Name, rule.Name, "matches", resourceName, "(", "match", match, ",", "exclude", exclude, ")")
				} else {
					fmt.Fprintln(out, "-", policy.Name, rule.Name, "doesn't match", resourceName, "(", "match", match, ",", "exclude", exclude, ")")
				}
			}
		}
	}
	fmt.Fprintln(out, "Done")
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
