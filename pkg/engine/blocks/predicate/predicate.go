package predicate

import (
	"github.com/kyverno/kyverno-json/pkg/engine"
)

type predicate[TREQUEST any, TRESPONSE any] struct {
	inner     engine.Engine[TREQUEST, TRESPONSE]
	predicate func(TREQUEST) bool
}

func (b *predicate[TREQUEST, TRESPONSE]) Run(request TREQUEST) []TRESPONSE {
	if !b.predicate(request) {
		return nil
	}
	return b.inner.Run(request)
}

func New[TREQUEST any, TRESPONSE any](inner engine.Engine[TREQUEST, TRESPONSE], condition func(TREQUEST) bool) engine.Engine[TREQUEST, TRESPONSE] {
	return &predicate[TREQUEST, TRESPONSE]{
		inner:     inner,
		predicate: condition,
	}
}
