//go:build js && wasm

package server

import (
	"context"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func RunWasm(_ context.Context, s Server) {
	wasmhttp.Serve(s.Handler())
}
