package null

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
)

type null[TREQUEST any, TRESPONSE any] struct{}

func (b *null[TREQUEST, TRESPONSE]) Run(_ context.Context, _ TREQUEST) []TRESPONSE {
	return nil
}

func New[TREQUEST any, TRESPONSE any]() engine.Engine[TREQUEST, TRESPONSE] {
	return &null[TREQUEST, TRESPONSE]{}
}
