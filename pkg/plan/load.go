package plan

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load(path string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var plan map[string]interface{}
	if err := json.Unmarshal(content, &plan); err != nil {
		return nil, err
	}
	return plan, nil
}
