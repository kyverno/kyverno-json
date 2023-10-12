package playground

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	Policy   *v1alpha1.Policy `json:"policy"`
	Rule     v1alpha1.Rule    `json:"rule"`
	Resource interface{}      `json:"resource"`
	Failure  error            `json:"failure"`
	Error    error            `json:"error"`
}

func makeResponse(responses ...jsonengine.JsonEngineResponse) *Response {
	var response Response
	for _, result := range responses {
		response.Results = append(response.Results, Result(result))
	}
	return &response
}
