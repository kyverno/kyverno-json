package test

import (
	"os"
	"path/filepath"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"gopkg.in/yaml.v2"
)

func loadTests(paths []string, fileName string) (TestCases, error) {
	var tests []TestCase
	for _, path := range paths {
		t, err := loadLocalTest(filepath.Clean(path), fileName)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t...)
	}
	return tests, nil
}

func loadLocalTest(path string, fileName string) (TestCases, error) {
	var tests []TestCase
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			ps, err := loadLocalTest(filepath.Join(path, file.Name()), fileName)
			if err != nil {
				return nil, err
			}
			tests = append(tests, ps...)
		} else if file.Name() == fileName {
			tests = append(tests, loadTest(filepath.Join(path, fileName), path))
		}
	}
	return tests, nil
}

func loadTest(path, currDir string) TestCase {
	var yamlBytes []byte
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return TestCase{
			Path: path,
			Err:  err,
		}
	}
	yamlBytes = data
	var test v1alpha1.Tests
	if err := yaml.Unmarshal(yamlBytes, &test); err != nil {
		return TestCase{
			Path: path,
			Err:  err,
		}
	}

	// Update the paths of policies and payloads based on current path
	for i := range test.Tests {
		for j := range test.Tests[i].Policies {
			test.Tests[i].Policies[j] = filepath.Join(currDir, test.Tests[i].Policies[j])
		}
		test.Tests[i].Payload = filepath.Join(currDir, test.Tests[i].Payload)
	}

	return TestCase{
		Path:  path,
		Tests: &test,
	}
}
