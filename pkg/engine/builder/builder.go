package builder

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/constant"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/function"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/predicate"
)

type Engine[TREQUEST any, TRESPONSE any] struct {
	engine.Engine[TREQUEST, TRESPONSE]
}

func new[TREQUEST any, TRESPONSE any](engine engine.Engine[TREQUEST, TRESPONSE]) Engine[TREQUEST, TRESPONSE] {
	return Engine[TREQUEST, TRESPONSE]{engine}
}

func Constant[TREQUEST any, TRESPONSE any](responses ...TRESPONSE) Engine[TREQUEST, TRESPONSE] {
	return new(constant.New[TREQUEST](responses...))
}

func (inner Engine[TREQUEST, TRESPONSE]) Predicate(condition func(context.Context, TREQUEST) bool) Engine[TREQUEST, TRESPONSE] {
	return new(predicate.New(inner, condition))
}

func Function[TREQUEST any, TRESPONSE any](f func(context.Context, TREQUEST) TRESPONSE) Engine[TREQUEST, TRESPONSE] {
	return new(function.New(f))
}
