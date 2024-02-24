package scan

import (
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

type Report struct {
	Resource any            `json:"resource"`
	Results  []PolicyReport `json:"results,omitempty"`
}

type PolicyReport struct {
	Policy *v1alpha1.ValidatingPolicy `json:"policy"`
	Rules  []RuleReport               `json:"rules,omitempty"`
}

type RuleReport struct {
	Rule       v1alpha1.ValidatingRule `json:"rule"`
	Identifier string                  `json:"identifier,omitempty"`
	Error      string                  `json:"error,omitempty"`
	Violations []ViolationReport       `json:"violations,omitempty"`
}

type ViolationReport struct {
	Message string        `json:"message,omitempty"`
	Errors  []ErrorReport `json:"errors,omitempty"`
}

type ErrorReport struct {
	Type   string `json:"type,omitempty"`
	Field  string `json:"field,omitempty"`
	Value  any    `json:"value,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func ToReport(response jsonengine.Response) Report {
	report := Report{
		Resource: response.Resource,
	}
	for _, policy := range response.Policies {
		policyReport := PolicyReport{
			Policy: policy.Policy,
		}
		for _, rule := range policy.Rules {
			ruleReport := RuleReport{
				Rule:       rule.Rule,
				Identifier: rule.Identifier,
			}
			if rule.Error != nil {
				ruleReport.Error = rule.Error.Error()
			}
			for _, violation := range rule.Violations {
				violationReport := ViolationReport{
					Message: violation.Message,
				}
				for _, err := range violation.ErrorList {
					if err != nil {
						violationReport.Errors = append(violationReport.Errors, ErrorReport{
							Type:   string(err.Type),
							Field:  err.Field,
							Value:  err.BadValue,
							Detail: err.Detail,
						})
					}
				}
				ruleReport.Violations = append(ruleReport.Violations, violationReport)
			}
			policyReport.Rules = append(policyReport.Rules, ruleReport)
		}
		report.Results = append(report.Results, policyReport)
	}
	return report
}

func ToReports(responses ...jsonengine.Response) (reports []Report) {
	for _, response := range responses {
		reports = append(reports, ToReport(response))
	}
	return reports
}
