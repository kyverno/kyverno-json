package template

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
)

var (
	funcs         = GetFunctions(context.Background())
	defaultCaller = interpreter.NewFunctionCaller(funcs...)
)

type Option func(options) options

type options struct {
	functionCaller interpreter.FunctionCaller
	engine         string
}

func WithFunctionCaller(functionCaller interpreter.FunctionCaller) Option {
	return func(o options) options {
		o.functionCaller = functionCaller
		return o
	}
}

func WithEngine(engine string) Option {
	return func(o options) options {
		o.engine = engine
		return o
	}
}

func buildOptions(opts ...Option) options {
	var o options
	for _, opt := range opts {
		if opt != nil {
			o = opt(o)
		}
	}
	if o.functionCaller == nil {
		o.functionCaller = defaultCaller
	}
	if o.engine == "" {
		o.engine = "jp"
	}
	return o
}
