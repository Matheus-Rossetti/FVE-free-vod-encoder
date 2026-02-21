package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func videoHandler(jobQueue chan<- core.VideoJob) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var request struct {
			VideoUrl string `json:"video_url"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		}

		// Download video and store in downloads/
		// then get the absolutepath for it and add it to
		// the queue

		jobQueue <- core.VideoJob{
			Source:            "REST",
			AbsoluteVideoPath: "",
		}

		fmt.Fprintf(w, "Video {video.name} added to internal queue and will be processed soon!")
	})
}
