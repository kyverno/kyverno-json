package constant

import (
	"github.com/eddycharly/tf-kyverno/pkg/engine"
)

type constant[TREQUEST any, TRESPONSE any] struct {
	responses []TRESPONSE
}

func (b *constant[TREQUEST, TRESPONSE]) Run(_ TREQUEST) []TRESPONSE {
	return b.responses
}

func New[TREQUEST any, TRESPONSE any](responses ...TRESPONSE) engine.Engine[TREQUEST, TRESPONSE] {
	return &constant[TREQUEST, TRESPONSE]{
		responses: responses,
	}
}
