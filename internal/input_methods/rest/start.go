package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job, port string) {
	fmt.Printf("\nStarting REST input method at http://localhost:%v\n", port)

	mux := http.NewServeMux()
	AddRoutes(mux, downloadQueue)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go server.ListenAndServe() // Run server

	<-ctx.Done() // Waits for term or int signal
	shutdownServer(server)
}
