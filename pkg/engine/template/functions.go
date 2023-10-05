package template

import (
	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template/functions"
)

func GetFunctions() []jpfunctions.FunctionEntry {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, jpfunctions.GetDefaultFunctions()...)
	funcs = append(funcs, functions.GetFunctions()...)
	return funcs
}
