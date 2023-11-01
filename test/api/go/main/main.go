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
    match:
      all:
        (input.method): "DELETE"
    assert:
      all:
      - check:
          role: "admin"
`

func main() {
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

	// create a JsonEngineRequest
	request := jsonengine.JsonEngineRequest{
		Resources: []interface{}{payload},
		Policies:  policies,
	}

	// create a J
	engine := jsonengine.New()

	responses := engine.Run(context.Background(), request)

	logger := log.Default()
	for _, resp := range responses {
		if resp.Error != nil {
			// ...handle execution error
			logger.Printf("policy error: %v", resp.Error)
		}

		if resp.Failure != nil {
			// ...handle policy failure
			logger.Printf("policy failure: %v", resp.Failure)
		}
	}
}
