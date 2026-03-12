package rest

import (
	"fmt"
	"net/http"
)

func (r *rest) Start() {
	addrs := fmt.Sprintf("http://localhost:%v", r.port)
	r.slog.Info("Starting... ", "Addrs", addrs)

	mux := http.NewServeMux()
	r.AddRoutes(mux)

	server := &http.Server{
		Addr:     r.port,
		Handler:  mux,
		ErrorLog: r.log,
	}

	go server.ListenAndServe() // Run server
	r.slog.Info("Waiting for requests!")

	<-r.ctx.Done() // Waits for term or int signal
	r.shutdownServer(server)
}
