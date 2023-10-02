package functions

import (
	"errors"
)

func at(arguments []interface{}) (interface{}, error) {
	sliceArg := arguments[0]
	indexArg := arguments[1]
	if slice, ok := sliceArg.([]interface{}); !ok {
		return nil, errors.New("invalid type, first argument must be an array")
	} else if index, ok := indexArg.(int); !ok {
		return nil, errors.New("invalid type, second argument must be an int")
	} else {
		return slice[index], nil
	}
}
