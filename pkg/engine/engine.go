package engine

// TODO:
// - tracing
// - explain
// - assertion tree

type Engine[TREQUEST any, TRESPONSE any] interface {
	Run(TREQUEST) []TRESPONSE
}
