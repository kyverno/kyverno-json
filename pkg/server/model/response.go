package model

import (
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

type Response struct {
	Results []RuleResponse `json:"results"`
}

type RuleResponse struct {
	PolicyName string                  `json:"policy"`
	RuleName   string                  `json:"rule"`
	Identifier string                  `json:"identifier,omitempty"`
	Result     jsonengine.PolicyResult `json:"result"`
	Message    string                  `json:"message"`
}

func MakeResponse(from ...jsonengine.Response) Response {
	var response Response
	for _, resource := range from {
		for _, policy := range resource.Policies {
			for _, rule := range policy.Rules {
				ruleResponse := RuleResponse{
					PolicyName: policy.Policy.Name,
					RuleName:   rule.Rule.Name,
					Identifier: rule.Identifier,
					Result:     makeResult(rule),
					Message:    makeMessage(rule),
				}
				response.Results = append(response.Results, ruleResponse)
			}
		}
	}
	return response
}

func makeResult(rule jsonengine.RuleResponse) jsonengine.PolicyResult {
	if rule.Error != nil {
		return jsonengine.StatusError
	}
	if len(rule.Violations) != 0 {
		return jsonengine.StatusFail
	}
	return jsonengine.StatusPass
}

func makeMessage(rule jsonengine.RuleResponse) string {
	if rule.Error != nil {
		return rule.Error.Error()
	}
	if len(rule.Violations) != 0 {
		return rule.Violations.Error()
	}
	return ""
}
