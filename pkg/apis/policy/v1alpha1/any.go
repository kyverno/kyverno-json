package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/json"
)

func deepCopy(in any) any {
	if in == nil {
		return nil
	}
	switch in := in.(type) {
	case string:
		return in
	case int:
		return in
	case int32:
		return in
	case int64:
		return in
	case float32:
		return in
	case float64:
		return in
	case bool:
		return in
	case []any:
		var out []any
		for _, in := range in {
			out = append(out, deepCopy(in))
		}
		return out
	case map[string]any:
		out := map[string]any{}
		for k, in := range in {
			out[k] = deepCopy(in)
		}
		return out
	}
	panic(fmt.Sprintf("deep copy failed - unrecognized type %T", in))
}

// Any can be any type.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
type Any struct {
	// +optional
	value any `json:"-"`
}

func NewAny(value any) Any {
	return Any{value}
}

func (t *Any) Value() any {
	return t.value
}

func (in *Any) DeepCopyInto(out *Any) {
	out.value = deepCopy(in.value)
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
