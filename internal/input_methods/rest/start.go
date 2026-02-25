package rest

import (
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(jobQueue chan<- core.Job, port string) {
	fmt.Printf("\nStarting REST input method in port %v...", port)

	AddRoutes(jobQueue)
	http.ListenAndServe(port, nil)
}
