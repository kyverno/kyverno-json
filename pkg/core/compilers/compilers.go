package compilers

import (
	"github.com/kyverno/kyverno-json/pkg/core/compilers/cel"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
)

var DefaultCompilers = Compilers{
	Jp:  jp.NewCompiler(),
	Cel: cel.NewCompiler(),
}

type Compilers struct {
	Jp  jp.Compiler
	Cel cel.Compiler
}

func (c Compilers) Compiler(compiler string) Compiler {
	switch compiler {
	case "":
		return nil
	case expression.CompilerJP:
		return c.Jp
	case expression.CompilerCEL:
		return c.Cel
	default:
		return c.Jp
	}
}
