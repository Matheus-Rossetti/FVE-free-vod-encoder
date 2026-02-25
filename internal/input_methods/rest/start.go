package rest

import (
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(downloadQueue chan<- core.DownloadJob, port string) {
	fmt.Printf("\nStarting REST input method in port %v...", port)

	AddRoutes(downloadQueue)
	http.ListenAndServe(port, nil)
}
