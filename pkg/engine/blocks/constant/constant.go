package constant

import (
	"context"

	"github.com/kyverno/kyverno-json/pkg/engine"
)

type constant[TREQUEST any, TRESPONSE any] struct {
	response TRESPONSE
}

func (b *constant[TREQUEST, TRESPONSE]) Run(_ context.Context, _ TREQUEST) TRESPONSE {
	return b.response
}

func New[TREQUEST any, TRESPONSE any](response TRESPONSE) engine.Engine[TREQUEST, TRESPONSE] {
	return &constant[TREQUEST, TRESPONSE]{
		response: response,
	}
}
