package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/json"
)

// Any can be any type.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
type Any struct {
	_value any
}

func NewAny(value any) Any {
	return Any{value}
}

func (t *Any) Value() any {
	return t._value
}

func (in *Any) DeepCopyInto(out *Any) {
	out._value = deepCopy(in._value)
}

func (in *Any) DeepCopy() *Any {
	if in == nil {
		return nil
	}
	out := new(Any)
	in.DeepCopyInto(out)
	return out
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
	return nil
}
