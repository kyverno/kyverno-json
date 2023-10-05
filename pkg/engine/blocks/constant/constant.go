package constant

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
)

type constant[TREQUEST any, TRESPONSE any] struct {
	responses []TRESPONSE
}

func (b *constant[TREQUEST, TRESPONSE]) Run(_ context.Context, _ TREQUEST) []TRESPONSE {
	return b.responses
}

func New[TREQUEST any, TRESPONSE any](responses ...TRESPONSE) engine.Engine[TREQUEST, TRESPONSE] {
	return &constant[TREQUEST, TRESPONSE]{
		responses: responses,
	}
}
