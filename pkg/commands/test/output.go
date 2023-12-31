package test

import (
	"fmt"
	"io"
	"strings"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/commands/test/output/color"
	"github.com/kyverno/kyverno-json/pkg/commands/test/output/table"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

type TestResponseOutput struct {
	Test      v1alpha1.Test
	Responses []jsonengine.RuleResponse
}

func printTestResults(out io.Writer, testResponse []TestResponseOutput, rc *resultCounts) (table.Table, error) {
	printer := table.NewTablePrinter(out)
	var resultsTable table.Table
	testCount := 1
	var rows []table.Row
	for _, tr := range testResponse {
		for _, r := range tr.Responses {
			success := (tr.Test.Result == v1alpha1.PolicyResult(r.Result))
			row := table.Row{
				ID:        testCount,
				Policy:    color.Policy("", r.PolicyName),
				Rule:      color.Rule(r.RuleName),
				Payload:   tr.Test.Payload,
				IsFailure: !success,
			}

			if success {
				row.Result = color.ResultPass()
				if tr.Test.Result == v1alpha1.StatusSkip {
					rc.Skip++
				} else {
					rc.Pass++
				}
			} else {
				row.Result = color.ResultFail()
				row.Reason = fmt.Sprintf("Expected: %s, recieved: %s %s", tr.Test.Result, r.Result, r.Message)
				rc.Fail++
			}
			testCount++
			rows = append(rows, row)
		}
		// if not found
		if len(rows) == 0 {
			row := table.Row{
				ID:        testCount,
				Policy:    color.Policy("", strings.Join(tr.Test.Policies, "")),
				Payload:   tr.Test.Payload,
				IsFailure: true,
				Result:    color.ResultFail(),
				Reason:    color.NotFound(),
			}
			testCount++
			resultsTable.Add(row)
			rc.Fail++
		}
	}
	resultsTable.Add(rows...)
	fmt.Fprintln(out)
	printer.Print(resultsTable.GetRows())
	fmt.Fprintln(out)
	return resultsTable, nil
}

func printFailedTestResult(out io.Writer, resultsTable table.Table) {
	printer := table.NewTablePrinter(out)
	for i := range resultsTable.Rows {
		resultsTable.Rows[i].ID = i + 1
	}
	fmt.Fprintf(out, "Aggregated Failed Test Cases : ")
	fmt.Fprintln(out)
	printer.Print(resultsTable.GetRows())
}
