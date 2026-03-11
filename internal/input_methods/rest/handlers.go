package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods"
)

func videoHandler(downloadQueue chan<- *core.Job) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var request struct {
			VideoUri     string `json:"video_url"`
			S3KeyStarter string `json:"s3_key_starter"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.VideoUri == "" {
			http.Error(w, "Package needs 'video_url'", http.StatusUnprocessableEntity)
			return
		}

		uriType := input_methods.CategorizeUri(request.VideoUri)
		if uriType == "unsupported uri" {
			fmt.Fprintf(w, "Unsupported URI Type\nAccepts: 'http', 'https', 'file' and '/' (local file)")
		}

		// BUILD JOB, GOTTA REFACTOR

		job := core.NewJob()

		job.DownloadJob.Source = "REST"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = request.VideoUri

		job.UploadJob.S3KeyStarter = request.S3KeyStarter

		downloadQueue <- job

		fmt.Fprintf(w, "Video %v added to internal queue and will be processed soon!", request.VideoUri)
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}
