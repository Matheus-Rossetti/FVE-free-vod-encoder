package rest

import (
	"fmt"
	"net/http"
	"time"
)

func (r *rest) Start() {
	r.slog.Info("Starting... ", "Addrs", fmt.Sprintf("http://localhost:%v", r.port))

	mux := http.NewServeMux()
	r.AddRoutes(mux)

	server := &http.Server{
		Addr:     r.port,
		Handler:  mux,
		ErrorLog: r.log,
	}

	go server.ListenAndServe()

	r.slog.Info("Waiting for requests!")

	<-r.ctx.Done()
	time.Sleep(time.Second * 2)
	r.shutdownServer(server) // already logs the shutdown
}
