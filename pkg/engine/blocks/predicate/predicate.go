package predicate

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
)

type predicate[TREQUEST any, TRESPONSE any] struct {
	inner     engine.Engine[TREQUEST, TRESPONSE]
	predicate func(context.Context, TREQUEST) bool
}

func (b *predicate[TREQUEST, TRESPONSE]) Run(ctx context.Context, request TREQUEST) TRESPONSE {
	if !b.predicate(ctx, request) {
		var none TRESPONSE
		return none
	}
	return b.inner.Run(ctx, request)
}

func New[TREQUEST any, TRESPONSE any](inner engine.Engine[TREQUEST, TRESPONSE], condition func(context.Context, TREQUEST) bool) engine.Engine[TREQUEST, TRESPONSE] {
	return &predicate[TREQUEST, TRESPONSE]{
		inner:     inner,
		predicate: condition,
	}
}
