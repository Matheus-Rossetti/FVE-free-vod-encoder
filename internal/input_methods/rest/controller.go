package rest

import (
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func AddRoutes(downloadQueue chan<- core.DownloadJob) {
	http.HandleFunc("POST /video", videoHandler(downloadQueue))
}
