package function

import (
	"github.com/kyverno/kyverno-json/pkg/engine"
)

type function[TREQUEST any, TRESPONSE any] struct {
	function func(TREQUEST) TRESPONSE
}

func (b *function[TREQUEST, TRESPONSE]) Run(request TREQUEST) []TRESPONSE {
	return []TRESPONSE{b.function(request)}
}

func New[TREQUEST any, TRESPONSE any](f func(TREQUEST) TRESPONSE) engine.Engine[TREQUEST, TRESPONSE] {
	return &function[TREQUEST, TRESPONSE]{
		function: f,
	}
}
