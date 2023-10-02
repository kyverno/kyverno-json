package functions

import (
	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

func GetFunctions() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: "at",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpArray}},
			// TODO: we should introduce a JpInteger type
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler: at,
	}}
}
