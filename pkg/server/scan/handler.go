package scan

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/server/model"
	"github.com/loopfz/gadgeto/tonic"
)

func newHandler(policyProvider PolicyProvider) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *Request) (*model.Response, error) {
		// check input
		if in == nil {
			return nil, errors.New("input is null")
		}
		if in.Payload == nil {
			return nil, errors.New("input payload is null")
		}
		payload := in.Payload
		// apply pre processors
		for _, preprocessor := range in.Preprocessors {
			result, err := template.Execute(context.Background(), preprocessor, payload, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to execute prepocessor (%s) - %w", preprocessor, err)
			}
			if result == nil {
				return nil, fmt.Errorf("prepocessor resulted in `null` payload (%s)", preprocessor)
			}
			payload = result
		}
		// load resources
		var resources []any
		if slice, ok := payload.([]any); ok {
			resources = slice
		} else {
			resources = append(resources, payload)
		}
		// load policies
		policies, err := policyProvider.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get policies (%w)", err)
		}
		var pols []*v1alpha1.ValidatingPolicy
		for i := range policies {
			pols = append(pols, &policies[i])
		}
		// run engine
		e := jsonengine.New()
		var results []jsonengine.Response
		for _, resource := range resources {
			results = append(results, e.Run(context.Background(), jsonengine.Request{
				Resource: resource,
				Policies: pols,
			}))
		}
		// TODO: return HTTP 403 for policy failure and HTTP 406 for policy errors
		response := model.MakeResponse(results...)
		return &response, nil
	}, http.StatusOK), nil
}
