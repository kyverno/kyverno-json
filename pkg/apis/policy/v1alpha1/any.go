package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/core/projection"
	hashutils "github.com/kyverno/kyverno-json/pkg/utils/hash"
	"k8s.io/apimachinery/pkg/util/json"
)

// Any can be any type.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
type Any struct {
	_value any
	_hash  string
}

func NewAny(value any) Any {
	return Any{
		_value: value,
		_hash:  hashutils.Hash(value),
	}
}

func (t *Any) Compile(compiler func(string, any, string) (projection.ScalarHandler, error), defaultCompiler string) (projection.ScalarHandler, error) {
	return compiler(t._hash, t._value, defaultCompiler)
}

func (a *Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a._value)
}

func (a *Any) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a._value = v
	a._hash = hashutils.Hash(a._value)
	return nil
}

func (in *Any) DeepCopyInto(out *Any) {
	out._value = deepCopy(in._value)
	out._hash = in._hash
}

func (in *Any) DeepCopy() *Any {
	if in == nil {
		return nil
	}
	out := new(Any)
	in.DeepCopyInto(out)
	return out
}
