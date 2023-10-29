package scan

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/loopfz/gadgeto/tonic"
)

func newHandler(policyProvider PolicyProvider) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *Request) (*Response, error) {
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
		var resources []interface{}
		if slice, ok := payload.([]interface{}); ok {
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
		results := e.Run(context.Background(), jsonengine.JsonEngineRequest{
			Resources: resources,
			Policies:  pols,
		})
		return makeResponse(results...), nil
	}, http.StatusOK), nil
}
