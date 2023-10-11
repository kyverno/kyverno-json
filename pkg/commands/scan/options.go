package scan

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/payload"
	"github.com/kyverno/kyverno-json/pkg/policy"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/output/pluralize"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

type options struct {
	payload       string
	preprocessors []string
	policies      []string
	selectors     []string
}

func (c *options) run(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Loading policies ...")
	policies, err := policy.Load(c.policies...)
	if err != nil {
		return err
	}
	selector := labels.Everything()
	if len(c.selectors) != 0 {
		parsed, err := labels.Parse(strings.Join(c.selectors, ","))
		if err != nil {
			return err
		}
		selector = parsed
	}
	{
		var filteredPolicies []*v1alpha1.Policy
		for _, policy := range policies {
			if selector.Matches(labels.Set(policy.Labels)) {
				filteredPolicies = append(filteredPolicies, policy)
			}
		}
		policies = filteredPolicies
	}
	fmt.Fprintln(out, "Loading payload ...")
	payload, err := payload.Load(c.payload)
	if err != nil {
		return err
	}
	if payload == nil {
		return errors.New("payload is `null`")
	}
	fmt.Fprintln(out, "Pre processing ...")
	for _, preprocessor := range c.preprocessors {
		result, err := template.Execute(context.Background(), preprocessor, payload, nil)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("prepocessor resulted in `null` payload (%s)", preprocessor)
		}
		payload = result
	}
	var resources []interface{}
	if slice, ok := payload.([]interface{}); ok {
		resources = slice
	} else {
		resources = append(resources, payload)
	}
	fmt.Fprintln(out, "Running", "(", "evaluating", len(resources), pluralize.Pluralize(len(resources), "resource", "resources"), "against", len(policies), pluralize.Pluralize(len(policies), "policy", "policies"), ")", "...")
	e := jsonengine.New()
	responses := e.Run(context.Background(), jsonengine.JsonEngineRequest{
		Resources: resources,
		Policies:  policies,
	})
	for _, response := range responses {
		resourceName, _, _ := unstructured.NestedString(response.Resource.(map[string]interface{}), "address")
		if response.Failure != nil {
			fmt.Fprintln(out, "-", response.Policy.Name, "/", response.Rule.Name, "/", resourceName, "ERROR:", response.Failure)
		} else {
			if response.Error == nil {
				fmt.Fprintln(out, "-", response.Policy.Name, "/", response.Rule.Name, "/", resourceName, "PASSED")
			} else {
				fmt.Fprintln(out, "-", response.Policy.Name, "/", response.Rule.Name, "/", resourceName, "FAILED:", response.Error)
			}
		}
	}
	fmt.Fprintln(out, "Done")
	return nil
}
