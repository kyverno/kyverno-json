package scan

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
)

type PolicyProvider interface {
	Get() ([]v1alpha1.Policy, error)
}
