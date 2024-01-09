package main

import (
	"context"
	"encoding/json"
	"log"

	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/policy"
)

const policyYAML = `
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: authz
spec:
  rules:
  - name: delete-checks
    identifier: "name"
    match:
      all:
        (input.method): "DELETE"
    assert:
      all:
      - check:
          role: "admin"
`

func main() {
	// load policies
	policies, err := policy.Parse([]byte(policyYAML))
	if err != nil {
		panic(err)
	}

	// load payloads
	requestJSON := `{
		"name": "Annie",
		"role": "admin",
		"input": {
			"method": "DELETE",
			"path":   "/red-files"
		}
	}`

	var payload interface{}
	if err := json.Unmarshal([]byte(requestJSON), &payload); err != nil {
		panic(err)
	}

	// create a Request
	request := jsonengine.Request{
		Resource: payload,
		Policies: policies,
	}

	// create an engine
	engine := jsonengine.New()

	// apply polices to get the response
	responses := engine.Run(context.Background(), request)

	// process the engine response
	logger := log.Default()
	for _, r := range responses {
		if r.Result == jsonengine.StatusFail {
			logger.Printf("fail: %s/%s -> %s: %s", r.PolicyName, r.RuleName, r.Identifier, r.Message)
		} else if r.Result == jsonengine.StatusError {
			logger.Printf("error: %s/%s -> %s: %s", r.PolicyName, r.RuleName, r.Identifier, r.Message)
		} else {
			logger.Printf("%s: %s/%s -> %s", r.Result, r.PolicyName, r.RuleName, r.Identifier)
		}
	}
}
