// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm

package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	server "github.com/kyverno/kyverno-json/pkg/server/wasm"
)

func main() {
	// initialise gin framework
	// gin.SetMode(c.ginFlags.mode)
	// tonic.SetBindHook(tonic.DefaultBindingHookMaxBodyBytes(int64(c.ginFlags.maxBodySize)))
	// create server
	server, err := server.New(true, true)
	if err != nil {
		panic(err)
	}
	// register playground routes
	if err := server.AddPlaygroundRoutes(); err != nil {
		panic(err)
	}
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdown := server.Run(ctx)
	<-ctx.Done()
	stop()
	if shutdown != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}
}
