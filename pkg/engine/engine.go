package engine

type Engine[TREQUEST any, TRESPONSE any] interface {
	Run(TREQUEST) []TRESPONSE
}
