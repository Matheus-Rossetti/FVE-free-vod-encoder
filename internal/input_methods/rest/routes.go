package rest

import (
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func AddRoutes(mux *http.ServeMux, downloadQueue chan<- *core.Job) {
	mux.HandleFunc("GET /", checkHealth)
	mux.HandleFunc("POST /encode", videoHandler(downloadQueue))
}
