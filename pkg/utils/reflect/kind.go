package reflect

import (
	"reflect"
)

func GetKind(value interface{}) reflect.Kind {
	if value == nil {
		return reflect.Invalid
	}
	return reflect.TypeOf(value).Kind()
}
