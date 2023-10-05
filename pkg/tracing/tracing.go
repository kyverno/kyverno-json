package tracing

import (
	"context"
	"fmt"
)

type Tracer interface {
	Trace(string)
}

type contextKey struct{}

func WithTracer(ctx context.Context, tracer Tracer) context.Context {
	return context.WithValue(ctx, contextKey{}, tracer)
}

func Trace(ctx context.Context, a ...any) {
	value := ctx.Value(contextKey{})
	if value != nil {
		if tracer, ok := value.(Tracer); ok {
			tracer.Trace(fmt.Sprint(a...))
		}
	}
}

func Tracef(ctx context.Context, format string, a ...any) {
	value := ctx.Value(contextKey{})
	if value != nil {
		if tracer, ok := value.(Tracer); ok {
			tracer.Trace(fmt.Sprintf(format, a...))
		}
	}
}
