package functions

import (
	"errors"
	"fmt"
)

func jpfAt(arguments []any) (any, error) {
	if slice, ok := arguments[0].([]any); !ok {
		return nil, errors.New("invalid type, first argument must be an array")
	} else if index, ok := arguments[1].(float64); !ok {
		return nil, errors.New("invalid type, second argument must be an int")
	} else {
		i := int(index)
		if i < 0 || i >= len(slice) {
			return nil, fmt.Errorf("index out of range: %d (array length %d)", i, len(slice))
		}
		return slice[i], nil
	}
}
