package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job, port string) {
	fmt.Printf("\nStarting REST input method in port %v\n", port)

	AddRoutes(downloadQueue)
	http.ListenAndServe(port, nil)
}
