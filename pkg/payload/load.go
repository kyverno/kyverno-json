package payload

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load(path string) (interface{}, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var payload interface{}
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}
