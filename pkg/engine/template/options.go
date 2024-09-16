package template

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
)

var (
	funcs         = GetFunctions(context.Background())
	defaultCaller = interpreter.NewFunctionCaller(funcs...)
)

type Option func(Options) Options

type Options struct {
	FunctionCaller interpreter.FunctionCaller
}

func WithFunctionCaller(functionCaller interpreter.FunctionCaller) Option {
	return func(o Options) Options {
		o.FunctionCaller = functionCaller
		return o
	}
}

func BuildOptions(opts ...Option) Options {
	var o Options
	for _, opt := range opts {
		if opt != nil {
			o = opt(o)
		}
	}
	if o.FunctionCaller == nil {
		o.FunctionCaller = defaultCaller
	}
	return o
}
