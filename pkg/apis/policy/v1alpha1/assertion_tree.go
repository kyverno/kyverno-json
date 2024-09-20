package v1alpha1

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"

	"github.com/kyverno/kyverno-json/pkg/core/assertion"
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

func hash(in any) string {
	if in == nil {
		return ""
	}
	bytes, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	hash := md5.Sum(bytes) //nolint:gosec
	return hex.EncodeToString(hash[:])
}

func NewAssertionTree(value any) AssertionTree {
	return AssertionTree{
		_tree: value,
		_hash: hash(value),
	}
}

func (t *AssertionTree) Compile(compiler func(string, any) (assertion.Assertion, error)) (assertion.Assertion, error) {
	return compiler(t._hash, t._tree)
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
	a._hash = hash(a._tree)
	return nil
}

func (in *AssertionTree) DeepCopyInto(out *AssertionTree) {
	out._tree = deepCopy(in._tree)
	out._hash = in._hash
}
