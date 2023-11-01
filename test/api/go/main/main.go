package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"gopkg.in/yaml.v2"
)

func main() {

	// load policies
	policyYAML := `
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

	var policy v1alpha1.ValidatingPolicy
	yaml.Unmarshal([]byte(policyYAML), &policy)

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
	json.Unmarshal([]byte(requestJSON), &payload)

	// create a JsonEngineRequest
	request := jsonengine.JsonEngineRequest{
		Resources: []interface{}{payload},
		Policies:  []*v1alpha1.ValidatingPolicy{&policy},
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
