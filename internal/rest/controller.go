package rest

import (
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func AddRoutes(jobQueue chan<- core.VideoJob) {
	http.HandleFunc("POST /video", videoHandler(jobQueue))
}
