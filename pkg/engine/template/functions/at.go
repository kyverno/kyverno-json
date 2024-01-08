package functions

import (
	"errors"
)

func jpfAt(arguments []any) (any, error) {
	if slice, ok := arguments[0].([]any); !ok {
		return nil, errors.New("invalid type, first argument must be an array")
	} else if index, ok := arguments[1].(int); !ok {
		return nil, errors.New("invalid type, second argument must be an int")
	} else {
		return slice[index], nil
	}
}
