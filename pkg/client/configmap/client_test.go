package configmap

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestConfigMapClient(t *testing.T) {
	// Create a fake clientset
	clientset := fake.NewSimpleClientset()

	// Create a test ConfigMap
	testConfigMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Add the ConfigMap to the fake clientset
	_, err := clientset.CoreV1().ConfigMaps("default").Create(context.Background(), testConfigMap, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test ConfigMap: %v", err)
	}

	// Create the ConfigMap client
	client := NewClient(clientset)

	// Test cases
	tests := []struct {
		name          string
		namespace     string
		configMapName string
		wantErr       bool
		checkData     bool
	}{
		{
			name:          "Get existing ConfigMap",
			namespace:     "default",
			configMapName: "test-configmap",
			wantErr:       false,
			checkData:     true,
		},
		{
			name:          "ConfigMap not found",
			namespace:     "default",
			configMapName: "non-existent",
			wantErr:       true,
			checkData:     false,
		},
		{
			name:          "Invalid namespace",
			namespace:     "",
			configMapName: "test-configmap",
			wantErr:       true,
			checkData:     false,
		},
		{
			name:          "Invalid name",
			namespace:     "default",
			configMapName: "",
			wantErr:       true,
			checkData:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm, err := client.Get(context.Background(), tt.configMapName, tt.namespace)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Don't check data if we expected an error
			if tt.wantErr {
				return
			}

			// Verify ConfigMap data
			if tt.checkData {
				if cm == nil {
					t.Fatal("Expected ConfigMap to be non-nil")
				}

				if cm.Name != tt.configMapName {
					t.Errorf("Expected ConfigMap name %s, got %s", tt.configMapName, cm.Name)
				}

				if cm.Namespace != tt.namespace {
					t.Errorf("Expected ConfigMap namespace %s, got %s", tt.namespace, cm.Namespace)
				}

				// Check data entries
				if val, ok := cm.Data["key1"]; !ok || val != "value1" {
					t.Errorf("Expected cm.Data[\"key1\"] = \"value1\", got %s", val)
				}

				if val, ok := cm.Data["key2"]; !ok || val != "value2" {
					t.Errorf("Expected cm.Data[\"key2\"] = \"value2\", got %s", val)
				}
			}
		})
	}
}

// TestMockConfigMapClient creates a mock client for testing
func TestMockConfigMapClient(t *testing.T) {
	mockClient := &mockClient{}

	configMap, err := mockClient.Get(context.Background(), "any-name", "any-namespace")
	if err != nil {
		t.Fatalf("Expected no error from mock client, got: %v", err)
	}

	if configMap == nil {
		t.Fatal("Expected ConfigMap to be non-nil")
	}

	// Check that the mock returned the expected values
	if configMap.Name != "any-name" {
		t.Errorf("Expected ConfigMap name to be \"any-name\", got %s", configMap.Name)
	}

	if configMap.Namespace != "any-namespace" {
		t.Errorf("Expected ConfigMap namespace to be \"any-namespace\", got %s", configMap.Namespace)
	}

	// Check data
	if _, ok := configMap.Data["test-key"]; !ok {
		t.Error("Expected configMap.Data to contain \"test-key\"")
	}
}

// mockClient is a simple mock implementation for testing
type mockClient struct{}

// Get implements the Client interface
func (m *mockClient) Get(ctx context.Context, name, namespace string) (*v1.ConfigMap, error) {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"test-key": "test-value",
		},
	}, nil
}
