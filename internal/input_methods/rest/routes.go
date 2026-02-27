package rest

import (
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func AddRoutes(downloadQueue chan<- *core.Job) {
	http.HandleFunc("POST /encode", videoHandler(downloadQueue))
}
