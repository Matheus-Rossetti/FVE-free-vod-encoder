package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (r *rest) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", r.checkHealth)
	mux.HandleFunc("POST /encode", r.videoHandler)

	server := &http.Server{
		Addr:     r.port,
		Handler:  mux,
		ErrorLog: r.log,
	}
	go server.ListenAndServe()

	addr := fmt.Sprintf("http://localhost:%v", r.port)
	r.slog.Info("Waiting for requests!", "addr", addr)

	<-r.ctx.Done()
	time.Sleep(time.Second * 2)
	r.slog.Info("Shutting down...")
	server.Shutdown(context.Background())
}
