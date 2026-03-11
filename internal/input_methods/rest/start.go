package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job, port string) {
	fmt.Printf("\nStarting REST input method at http://localhost:%v\n", port)

	mux := http.NewServeMux()
	AddRoutes(mux, downloadQueue)

	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	go server.ListenAndServe() // Start server

	<-ctx.Done() // Waits for term or int signal
	fmt.Printf("\nREST input method shutting down...")

	serverCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := server.Shutdown(serverCtx)
	if err != nil {
		fmt.Printf("\nError shutting down REST input method: %v", err)
		return
	}

}
