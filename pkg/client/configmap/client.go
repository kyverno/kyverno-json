package configmap

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Client provides functions to interact with ConfigMaps
type Client interface {
	// Get fetches a ConfigMap from the Kubernetes API server
	Get(ctx context.Context, name, namespace string) (*v1.ConfigMap, error)
}

// client implements the Client interface
type client struct {
	clientset kubernetes.Interface
}

// NewClient creates a new instance of the ConfigMap client
func NewClient(clientset kubernetes.Interface) Client {
	return &client{
		clientset: clientset,
	}
}

// Get retrieves a ConfigMap from the API server
func (c *client) Get(ctx context.Context, name, namespace string) (*v1.ConfigMap, error) {
	if name == "" {
		return nil, fmt.Errorf("ConfigMap name cannot be empty")
	}
	if namespace == "" {
		return nil, fmt.Errorf("ConfigMap namespace cannot be empty")
	}

	return c.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}
