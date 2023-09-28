package plan

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Plan struct {
	Plan      map[string]interface{}
	Resources []interface{}
}

func Load(path string) (*Plan, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var plan map[string]interface{}
	if err := json.Unmarshal(content, &plan); err != nil {
		return nil, err
	}
	resources, ok, err := unstructured.NestedSlice(plan, "planned_values", "root_module", "resources")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to find resources in the plan")
	}
	return &Plan{
		Plan:      plan,
		Resources: resources,
	}, nil
}
