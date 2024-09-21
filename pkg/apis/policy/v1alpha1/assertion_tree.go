package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/core/assertion"
	hashutils "github.com/kyverno/kyverno-json/pkg/utils/hash"
	"k8s.io/apimachinery/pkg/util/json"
)

// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
// AssertionTree represents an assertion tree.
type AssertionTree struct {
	_tree any
	_hash string
}

func NewAssertionTree(value any) AssertionTree {
	return AssertionTree{
		_tree: value,
		_hash: hashutils.Hash(value),
	}
}

func (t *AssertionTree) Compile(compiler func(string, any, string) (assertion.Assertion, error), defaultCompiler string) (assertion.Assertion, error) {
	return compiler(t._hash, t._tree, defaultCompiler)
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
	a._hash = hashutils.Hash(a._tree)
	return nil
}

func (in *AssertionTree) DeepCopyInto(out *AssertionTree) {
	out._tree = deepCopy(in._tree)
	out._hash = in._hash
}
