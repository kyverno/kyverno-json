package payload

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eddycharly/json-kyverno/pkg/utils/file"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"gopkg.in/yaml.v3"
)

func Load(path string) (interface{}, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var payload interface{}
	switch {
	case file.IsJson(path):
		if err := json.Unmarshal(content, &payload); err != nil {
			return nil, err
		}
	case file.IsYaml(path):
		documents, err := yamlutils.SplitDocuments(content)
		if err != nil {
			return nil, err
		}
		var objects []interface{}
		for _, document := range documents {
			var object map[string]interface{}
			if err := yaml.Unmarshal(document, &object); err != nil {
				return nil, err
			}
			objects = append(objects, object)
		}
		if len(objects) == 1 {
			payload = objects[0]
		} else {
			payload = objects
		}
	default:
		return nil, fmt.Errorf("unrecognized payload format, must be yaml or json (%s)", path)
	}
	return payload, nil
}
