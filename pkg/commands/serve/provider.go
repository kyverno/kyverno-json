package serve

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type provider struct {
	client versioned.Interface
}

// TODO: use an informer/lister
func (p *provider) Get() ([]v1alpha1.Policy, error) {
	list, err := p.client.JsonV1alpha1().Policies().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
