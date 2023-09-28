package commands

import (
	"fmt"

	jsonengine "github.com/eddycharly/json-kyverno/pkg/json-engine"
	"github.com/eddycharly/json-kyverno/pkg/plan"
	"github.com/eddycharly/json-kyverno/pkg/policy"
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
	fmt.Fprintln(out, "-", len(plan.Resources), pluralize.Pluralize(len(plan.Resources), "resource", "resources"), "loaded")
	fmt.Fprintln(out, "Running ...")
	e := jsonengine.New()
	responses := e.Run(jsonengine.JsonEngineRequest{
		Plan:     plan,
		Policies: policies,
	})
	for _, response := range responses {
		resourceName, _, _ := unstructured.NestedString(response.Resource.(map[string]interface{}), "address")
		if response.Error == nil {
			fmt.Fprintln(out, "-", response.Policy.Name, "/", response.Rule.Name, "/", resourceName, "PASSED")
		} else {
			fmt.Fprintln(out, "-", response.Policy.Name, "/", response.Rule.Name, "/", resourceName, "FAILED:", response.Error)
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
