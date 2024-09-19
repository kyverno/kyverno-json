package v1alpha1

import (
	"context"
	"sync"

	"github.com/kyverno/kyverno-json/pkg/core/assertion"
	"k8s.io/apimachinery/pkg/util/json"
)

// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
// AssertionTree represents an assertion tree.
type AssertionTree struct {
	_tree      any
	_assertion func() (assertion.Assertion, error)
}

func NewAssertionTree(value any) AssertionTree {
	return AssertionTree{
		_tree: value,
		_assertion: sync.OnceValues(func() (assertion.Assertion, error) {
			return assertion.Parse(context.Background(), value)
		}),
	}
}

func (t *AssertionTree) Assertion() (assertion.Assertion, error) {
	if t._tree == nil {
		return nil, nil
	}
	return t._assertion()
}

func (a *AssertionTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(a._tree)
}

func (a *AssertionTree) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a._tree = v
	a._assertion = sync.OnceValues(func() (assertion.Assertion, error) {
		return assertion.Parse(context.Background(), v)
	})
	return nil
}

func (in *AssertionTree) DeepCopyInto(out *AssertionTree) {
	out._tree = deepCopy(in._tree)
	out._assertion = in._assertion
}
