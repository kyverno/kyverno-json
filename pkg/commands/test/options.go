package test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/commands/test/output/table"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/payload"
	"github.com/kyverno/kyverno-json/pkg/policy"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/output/color"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
)

type options struct {
	fileName    string
	removeColor bool
}

type resultCounts struct {
	Pass int
	Skip int
	Fail int
}

func (c *options) run(cmd *cobra.Command, dirPath []string) error {
	out := cmd.OutOrStdout()
	color.Init(c.removeColor)
	testcases, err := loadTests(dirPath, c.fileName)
	if err != nil {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Error loading tests:", err)
		return err
	}
	if len(testcases) == 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "No test yamls available")
		return nil
	}
	if errs := testcases.Errors(); len(errs) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Test errors:")
		for _, e := range errs {
			fmt.Fprintln(out, "  Path:", e.Path)
			fmt.Fprintln(out, "    Error:", e.Err)
		}
	}
	rc := &resultCounts{}
	var table table.Table
	e := jsonengine.New()
	for _, tc := range testcases {
		testRespOut := make([]TestResponseOutput, 0)
		for _, test := range tc.Tests.Tests {
			policies, err := policy.Load(test.Policies...)
			if err != nil {
				fmt.Println("Failed to load policy:", err.Error())
				return err
			}
			selector := labels.Everything()
			if len(test.Selectors) != 0 {
				parsed, err := labels.Parse(strings.Join(test.Selectors, ","))
				if err != nil {
					fmt.Println("Failed to parse labels:", err.Error())
					return err
				}
				selector = parsed
			}
			{
				var filteredPolicies []*v1alpha1.ValidatingPolicy
				for _, policy := range policies {
					if selector.Matches(labels.Set(policy.Labels)) {
						filteredPolicies = append(filteredPolicies, policy)
					}
				}
				policies = filteredPolicies
			}

			payload, err := payload.Load(test.Payload)
			if err != nil {
				fmt.Println("Failed to load payload:", err.Error())
				return err
			}
			if payload == nil {
				return errors.New("payload is `null`")
			}
			for _, preprocessor := range test.Preprocessors {
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

			responses := e.Run(context.Background(), jsonengine.Request{
				Resources: resources,
				Policies:  policies,
			})
			if len(responses) != 0 {
				testRespOut = append(testRespOut, TestResponseOutput{
					Test:      test,
					Responses: responses,
				})
			}
		}
		t, err := printTestResults(out, testRespOut, rc)
		if err != nil {
			return fmt.Errorf("failed to print test result (%w)", err)
		}
		table.AddFailed(t.Rows...)
	}
	fmt.Fprintf(out, "\nTest Summary: %d out of %d tests failed\n", rc.Fail, rc.Pass+rc.Skip+rc.Fail)
	fmt.Fprintln(out)
	if rc.Fail > 0 {
		printFailedTestResult(out, table)
		return fmt.Errorf("%d tests failed", rc.Fail)
	}
	return nil
}
