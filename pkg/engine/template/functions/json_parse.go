package functions

import (
	"encoding/json"
	"errors"
)

func jpfJsonParse(arguments []interface{}) (interface{}, error) {
	if data, ok := arguments[0].(string); !ok {
		return nil, errors.New("invalid type, first argument must be a string")
	} else {
		var result interface{}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
