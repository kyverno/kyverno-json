package playground

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/templating"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/server/model"
	"github.com/loopfz/gadgeto/tonic"
	"sigs.k8s.io/yaml"
)

func newHandler() (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *Request) (*model.Response, error) {
		// check input
		if in == nil {
			return nil, errors.New("input is null")
		}
		if in.Payload == "" {
			return nil, errors.New("input payload is null")
		}
		if in.Policy == "" {
			return nil, errors.New("input policy is null")
		}
		var payload any
		err := yaml.Unmarshal([]byte(in.Payload), &payload)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload (%w)", err)
		}
		// apply pre processors
		for _, preprocessor := range in.Preprocessors {
			result, err := templating.ExecuteJP(preprocessor, payload, nil, templating.NewCompiler(templating.CompilerOptions{}))
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
		// load policy
		var policy v1alpha1.ValidatingPolicy
		if err := yaml.Unmarshal([]byte(in.Policy), &policy); err != nil {
			return nil, fmt.Errorf("failed to parse policies (%w)", err)
		}
		// run engine
		e := jsonengine.New()
		var results []jsonengine.Response
		for _, resource := range resources {
			results = append(results, e.Run(context.Background(), jsonengine.Request{
				Resource: resource,
				Policies: []*v1alpha1.ValidatingPolicy{&policy},
			}))
		}
		response := model.MakeResponse(results...)
		return &response, nil
	}, http.StatusOK), nil
}
