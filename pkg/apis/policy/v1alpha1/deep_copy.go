package v1alpha1

import (
	"fmt"
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
