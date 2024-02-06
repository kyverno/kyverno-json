package scan

import (
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type PolicyProvider interface {
	Get() ([]v1alpha1.ValidatingPolicy, error)
}
