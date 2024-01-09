package engine

import (
	"context"
)

// TODO:
// - tracing
// - explain

type Engine[TREQUEST any, TRESPONSE any] interface {
	Run(context.Context, TREQUEST) TRESPONSE
}
