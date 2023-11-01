package scan

import (
	"net/http"

	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

type Response struct {
	Results []Result `json:"results"`
}

type PolicyResult string

type Result struct {
	PolicyName string       `json:"policy"`
	RuleName   string       `json:"rule"`
	Result     PolicyResult `json:"status"`
	Message    string       `json:"message"`
}

// Status specifies state of a policy result
const (
	StatusPass  PolicyResult = "pass"
	StatusFail  PolicyResult = "fail"
	StatusWarn  PolicyResult = "warn"
	StatusError PolicyResult = "error"
	StatusSkip  PolicyResult = "skip"
)

func makeResponse(responses ...jsonengine.JsonEngineResponse) (*Response, int) {
	var response Response
	failCount := 0
	errorCount := 0
	for _, r := range responses {
		status, msg := getStatusAndMessage(r)
		if status == StatusError {
			errorCount++
		} else if status == StatusFail {
			failCount++
		}

		response.Results = append(response.Results, Result{
			PolicyName: r.Policy.Name,
			RuleName:   r.Rule.Name,
			Result:     status,
			Message:    msg,
		})
	}

	httpStatus := http.StatusOK
	if failCount > 0 {
		httpStatus = http.StatusForbidden
	} else if errorCount > 0 {
		httpStatus = http.StatusNotAcceptable
	}

	return &response, httpStatus
}

func getStatusAndMessage(r jsonengine.JsonEngineResponse) (PolicyResult, string) {
	if r.Error != nil {
		return StatusError, r.Error.Error()
	}
	if r.Failure != nil {
		return StatusFail, r.Failure.Error()
	}
	return StatusPass, ""
}
