package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func videoHandler(downloadQueue chan<- core.DownloadJob) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var request struct {
			VideoUrl string `json:"video_url"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.VideoUrl == "" {
			http.Error(w, "Package needs 'video_url'", http.StatusUnprocessableEntity)
			return
		}

		downloadJob := core.NewDownloadJob(request.VideoUrl, "REST")
		downloadQueue <- downloadJob

		fmt.Fprintf(w, "Video {video.name} added to internal queue and will be processed soon!")
	})
}
