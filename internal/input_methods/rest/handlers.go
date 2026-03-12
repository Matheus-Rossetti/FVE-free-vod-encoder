package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (rest *rest) videoHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var request struct {
			VideoUri   string `json:"video_url"`
			KeyStarter string `json:"key_starter"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.VideoUri == "" || request.KeyStarter == "" {
			http.Error(w, "Needs 'video_url' and 'key_starter'", http.StatusUnprocessableEntity)
			return
		}

		uriType, err := core.CategorizeUri(request.VideoUri)
		if errors.Is(err, core.ErrUnsupportedUri) {
			http.Error(w, "Unsupported URI Type!\nAccepts: http, https and file.", http.StatusUnprocessableEntity)
		} else if err != nil {
			http.Error(w, "Error categorizing URI!\nAccepts: http, https and file.", http.StatusInternalServerError)
		}

		job := core.NewJob()

		job.DownloadJob.Source = "REST"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = request.VideoUri

		job.UploadJob.KeyStarter = request.KeyStarter

		rest.downloadQueue <- job

		fmt.Fprintf(w, "%v added to internal queue and will be processed soon!", request.VideoUri)
	})
}

func (rest *rest) checkHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}
