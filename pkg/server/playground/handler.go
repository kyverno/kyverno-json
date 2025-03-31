package playground

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/server/model"
	"github.com/loopfz/gadgeto/tonic"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Initialize a mock ConfigMapClient when the playground handler is created
func init() {
	// Use a mock implementation of the ConfigMapClient for playground
	// This way users can test policies that use ConfigMap references in the playground
	// without needing a real Kubernetes cluster
	jsonengine.SetConfigMapClient(&mockConfigMapClient{})
}

// mockConfigMapClient is a simple implementation of the ConfigMapClient interface
// that returns a predefined response for playground use
type mockConfigMapClient struct{}

// Get returns a mock ConfigMap for playground use
func (m *mockConfigMapClient) Get(ctx context.Context, name, namespace string) (*v1.ConfigMap, error) {
	// For playground purposes, return a mock ConfigMap with some sample data
	// This allows testing policies with ConfigMap references
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"key1":      "value1",
			"key2":      "value2",
			"allowList": "item1,item2,item3",
		},
	}, nil
}

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
			result, err := compilers.Execute(preprocessor, payload, nil, compilers.DefaultCompilers.Jp)
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
