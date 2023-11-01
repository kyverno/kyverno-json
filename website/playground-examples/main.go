package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Content struct {
	Policy  string `yaml:"policy"`
	Payload string `yaml:"payload"`
}

type Example struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Policy   string `json:"policy"`
	Payload  string `json:"payload"`
}

type Examples struct {
	Examples []Example `json:"examples"`
}

func load(file string) string {
	if filepath.Ext(file) == ".json" {
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		var obj interface{}
		if err := json.Unmarshal(content, &obj); err != nil {
			panic(err)
		}
		data, err := yaml.Marshal(obj)
		if err != nil {
			panic(err)
		}
		return string(data)
	} else if filepath.Ext(file) == ".yaml" || filepath.Ext(file) == ".yml" {
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		return string(content)
	}
	panic(fmt.Sprintf("unsupported file type %s", file))
}

func main() {
	bytes, err := os.ReadFile("playground-examples.yaml")
	if err != nil {
		panic(err)
	}
	var content map[string]map[string]Content
	if err := yaml.Unmarshal(bytes, &content); err != nil {
		panic(err)
	}
	var examples Examples
	for category, names := range content {
		for name, value := range names {
			examples.Examples = append(examples.Examples, Example{
				Category: category,
				Name:     name,
				Policy:   load(value.Policy),
				Payload:  load(value.Payload),
			})
		}
	}
	data, err := json.MarshalIndent(&examples, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("website/playground/assets/data.json", data, fs.ModePerm); err != nil {
		panic(err)
	}
}
