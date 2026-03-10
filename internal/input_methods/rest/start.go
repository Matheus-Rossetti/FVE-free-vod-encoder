package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job, port string) {
	fmt.Printf("\nStarting REST input method in port %v\n", port)

	mux := http.NewServeMux()
	AddRoutes(mux, downloadQueue)

	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	go server.ListenAndServe()

	<-ctx.Done()

	serverCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	server.Shutdown(serverCtx)
}
