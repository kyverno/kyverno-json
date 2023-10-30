//go:build !js && !wasm

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run(_ context.Context, s Server, host string, port int) Shutdown {
	address := fmt.Sprintf("%v:%v", host, port)
	srv := &http.Server{
		Addr:              address,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Default().Printf("listening to requests on %s:%d", host, port)
	return srv.Shutdown
}
