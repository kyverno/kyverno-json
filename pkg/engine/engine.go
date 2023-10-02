package engine

// TODO:
// - tracing
// - explain

type Engine[TREQUEST any, TRESPONSE any] interface {
	Run(TREQUEST) []TRESPONSE
}
