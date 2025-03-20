package payload

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func Test_Load(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "payload-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tempDir)

	// Test 1: Load JSON file
	jsonContent := `{"name": "test", "string_value": "hello"}`
	jsonPath := filepath.Join(tempDir, "test.json")
	err = os.WriteFile(jsonPath, []byte(jsonContent), 0644)
	assert.NilError(t, err)

	result, err := Load(jsonPath)
	assert.NilError(t, err)
	resultMap, ok := result.(map[string]interface{})
	assert.Assert(t, ok, "expected map[string]interface{} for JSON result")
	assert.Equal(t, "test", resultMap["name"])
	assert.Equal(t, "hello", resultMap["string_value"])

	// Test 2: Load single document YAML file
	yamlContent := `name: test
string_value: hello`
	yamlPath := filepath.Join(tempDir, "test.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	assert.NilError(t, err)

	result, err = Load(yamlPath)
	assert.NilError(t, err)
	resultMap, ok = result.(map[string]interface{})
	assert.Assert(t, ok, "expected map[string]interface{} for YAML result")
	assert.Equal(t, "test", resultMap["name"])
	assert.Equal(t, "hello", resultMap["string_value"])

	// Test 3: Load multi-document YAML file
	multiYAMLContent := `---
name: test1
string_value: hello
---
name: test2
string_value: world
`
	multiYAMLPath := filepath.Join(tempDir, "multi.yaml")
	err = os.WriteFile(multiYAMLPath, []byte(multiYAMLContent), 0644)
	assert.NilError(t, err)

	result, err = Load(multiYAMLPath)
	assert.NilError(t, err)
	resultArray, ok := result.([]interface{})
	assert.Assert(t, ok, "expected []interface{} for multi-doc YAML result")
	assert.Equal(t, 2, len(resultArray))

	doc1, ok := resultArray[0].(map[string]interface{})
	assert.Assert(t, ok, "expected map[string]interface{} for first YAML doc")
	assert.Equal(t, "test1", doc1["name"])
	assert.Equal(t, "hello", doc1["string_value"])

	doc2, ok := resultArray[1].(map[string]interface{})
	assert.Assert(t, ok, "expected map[string]interface{} for second YAML doc")
	assert.Equal(t, "test2", doc2["name"])
	assert.Equal(t, "world", doc2["string_value"])

	// Test 4: File not found error
	_, err = Load(filepath.Join(tempDir, "nonexistent.json"))
	assert.Assert(t, err != nil, "expected error for nonexistent file")

	// Test 5: Unrecognized format error
	txtPath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(txtPath, []byte("plain text"), 0644)
	assert.NilError(t, err)

	_, err = Load(txtPath)
	assert.Assert(t, err != nil, "expected error for unrecognized format")
	assert.Assert(t, err.Error() != "", "error message should not be empty")

	// Test 6: Invalid JSON format
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSONPath, []byte(`{"name": "test", invalid}`), 0644)
	assert.NilError(t, err)

	_, err = Load(invalidJSONPath)
	assert.Assert(t, err != nil, "expected error for invalid JSON")

	// Test 7: Invalid YAML format
	invalidYAMLPath := filepath.Join(tempDir, "invalid.yaml")
	err = os.WriteFile(invalidYAMLPath, []byte(`name: test
  value: - invalid`), 0644)
	assert.NilError(t, err)

	_, err = Load(invalidYAMLPath)
	assert.Assert(t, err != nil, "expected error for invalid YAML")
}
