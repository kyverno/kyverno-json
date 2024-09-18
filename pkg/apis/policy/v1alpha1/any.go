package v1alpha1

import (
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/util/json"
)

// Any can be any type.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
type Any struct {
	value any `json:"-"`
}

func (t *Any) Value() any {
	return t.value
}

func (in *Any) DeepCopyInto(out *Any) {
	if err := copier.CopyWithOption(out, in, copier.Option{DeepCopy: true}); err != nil {
		panic("deep copy failed")
	}
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
	return json.Marshal(a.value)
}

func (a *Any) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a.value = v
	return nil
}
