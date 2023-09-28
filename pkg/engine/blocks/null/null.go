package null

import (
	"github.com/eddycharly/json-kyverno/pkg/engine"
)

type null[TREQUEST any, TRESPONSE any] struct{}

func (b *null[TREQUEST, TRESPONSE]) Run(_ TREQUEST) []TRESPONSE {
	return nil
}

func New[TREQUEST any, TRESPONSE any]() engine.Engine[TREQUEST, TRESPONSE] {
	return &null[TREQUEST, TRESPONSE]{}
}
