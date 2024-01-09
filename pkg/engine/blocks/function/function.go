package function

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
)

type function[TREQUEST any, TRESPONSE any] struct {
	function func(context.Context, TREQUEST) TRESPONSE
}

func (b *function[TREQUEST, TRESPONSE]) Run(ctx context.Context, request TREQUEST) TRESPONSE {
	return b.function(ctx, request)
}

func New[TREQUEST any, TRESPONSE any](f func(context.Context, TREQUEST) TRESPONSE) engine.Engine[TREQUEST, TRESPONSE] {
	return &function[TREQUEST, TRESPONSE]{
		function: f,
	}
}
