package jsonengine

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// mockConfigMapClient is a mock client for testing
type mockConfigMapClient struct{}

// Get returns a mock ConfigMap for testing
func (m *mockConfigMapClient) Get(ctx context.Context, name, namespace string) (*v1.ConfigMap, error) {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"environment": "production",
			"maxReplicas": "5",
		},
	}, nil
}

// TestConfigMapContextEntry tests the ConfigMap context entry functionality
func TestConfigMapContextEntry(t *testing.T) {
	// Set the mock ConfigMap client
	SetConfigMapClient(&mockConfigMapClient{})

	// Create a context entry that references a ConfigMap
	contextEntry := v1alpha1.ContextEntry{
		Name: "testConfig",
		ConfigMap: &v1alpha1.ConfigMapReference{
			Name:      "test-configmap",
			Namespace: "default",
		},
	}

	// Create compilers and path for testing
	compilers := compilers.DefaultCompilers
	path := field.NewPath("test")

	// Test compiling the context entry
	contextHandler, err := CompileContextEntry(path, compilers, contextEntry)
	if err != nil {
		t.Fatalf("Failed to compile ConfigMap context entry: %v", err)
	}

	// Create bindings and sample resource
	bindings := binding.NewBindings()
	resource := map[string]interface{}{"sample": "data"}

	// Apply the context handler to update bindings
	updatedBindings := contextHandler(resource, bindings)

	// Create a simple JMESPath expression to test accessing the ConfigMap data
	jp := compilers.Jp

	// Test accessing the ConfigMap name
	nameExpr, compileErr := jp.Compile("$testConfig.metadata.name")
	if compileErr != nil {
		t.Fatalf("Failed to compile JMESPath expression: %v", compileErr)
	}

	nameValue, evalErr := nameExpr(resource, updatedBindings)
	if evalErr != nil {
		t.Fatalf("Failed to evaluate expression: %v", evalErr)
	}

	if nameValue != "test-configmap" {
		t.Errorf("Expected ConfigMap name to be 'test-configmap', got %v", nameValue)
	}

	// Test accessing ConfigMap data
	envExpr, compileErr := jp.Compile("$testConfig.data.environment")
	if compileErr != nil {
		t.Fatalf("Failed to compile JMESPath expression: %v", compileErr)
	}

	envValue, evalErr := envExpr(resource, updatedBindings)
	if evalErr != nil {
		t.Fatalf("Failed to evaluate expression: %v", evalErr)
	}

	if envValue != "production" {
		t.Errorf("Expected environment to be 'production', got %v", envValue)
	}
}

func TestCompileContextEntry_NoClient(t *testing.T) {
	// Reset the client to nil
	SetConfigMapClient(nil)

	// Create a context entry with ConfigMap reference
	contextEntry := v1alpha1.ContextEntry{
		Name: "config",
		ConfigMap: &v1alpha1.ConfigMapReference{
			Name:      "test-config",
			Namespace: "default",
		},
	}

	// Create field path and compilers
	path := field.NewPath("test")
	compilers := compilers.DefaultCompilers

	// Compile the context entry - should fail because client is nil
	_, err := CompileContextEntry(path, compilers, contextEntry)
	if err == nil {
		t.Fatal("Expected error when ConfigMap client is nil, got nil")
	}

	// Restore the mock client for other tests
	SetConfigMapClient(&mockConfigMapClient{})
}

// TestConfigMapContextEntry_NoClient tests that an error is returned when the ConfigMap client is nil
func TestConfigMapContextEntry_NoClient(t *testing.T) {
	// Reset the client to nil
	SetConfigMapClient(nil)

	// Create a context entry that references a ConfigMap
	contextEntry := v1alpha1.ContextEntry{
		Name: "testConfig",
		ConfigMap: &v1alpha1.ConfigMapReference{
			Name:      "test-configmap",
			Namespace: "default",
		},
	}

	// Create compilers and path for testing
	compilers := compilers.DefaultCompilers
	path := field.NewPath("test")

	// Test compiling the context entry - should fail
	_, err := CompileContextEntry(path, compilers, contextEntry)
	if err == nil {
		t.Fatal("Expected error when ConfigMap client is nil, got nil")
	}

	// Restore the mock client for other tests
	SetConfigMapClient(&mockConfigMapClient{})
}
