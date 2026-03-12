package rest

import (
	"net/http"
)

func (r *rest) Start() {
	r.slog.Info("Starting... ", "Addrs", "http://localhost:"+r.port)

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
