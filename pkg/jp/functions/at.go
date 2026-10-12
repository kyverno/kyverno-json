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
	} else if idx := int(index); idx < 0 || idx >= len(slice) {
		return nil, fmt.Errorf("index %d out of range for array of length %d", idx, len(slice))
	} else {
		return slice[idx], nil
	}
}
