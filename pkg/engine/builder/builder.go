package builder

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/constant"
	"github.com/kyverno/kyverno-json/pkg/engine/blocks/function"
)

type Engine[TREQUEST any, TRESPONSE any] struct {
	engine.Engine[TREQUEST, TRESPONSE]
}

func new[TREQUEST any, TRESPONSE any](engine engine.Engine[TREQUEST, TRESPONSE]) Engine[TREQUEST, TRESPONSE] {
	return Engine[TREQUEST, TRESPONSE]{engine}
}

func Constant[TREQUEST any, TRESPONSE any](response TRESPONSE) Engine[TREQUEST, TRESPONSE] {
	return new(constant.New[TREQUEST](response))
}

func Function[TREQUEST any, TRESPONSE any](f func(context.Context, TREQUEST) TRESPONSE) Engine[TREQUEST, TRESPONSE] {
	return new(function.New(f))
}
