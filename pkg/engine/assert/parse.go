package assert

import (
	"reflect"

	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

func Parse(assertion interface{}) Assertion {
	switch reflectutils.GetKind(assertion) {
	case reflect.Slice:
		node := sliceNode{}
		valueOf := reflect.ValueOf(assertion)
		for i := 0; i < valueOf.Len(); i++ {
			node = append(node, Parse(valueOf.Index(i).Interface()))
		}
		return node
	case reflect.Map:
		node := mapNode{}
		iter := reflect.ValueOf(assertion).MapRange()
		for iter.Next() {
			node[iter.Key().Interface()] = Parse(iter.Value().Interface())
		}
		return node
	default:
		return &scalarNode{rhs: assertion}
	}
}
