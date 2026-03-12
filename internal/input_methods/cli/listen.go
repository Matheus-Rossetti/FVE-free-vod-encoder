package cli

import (
	"fmt"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods"
)

func (c *cli) ListenForInput() {
	c.slog.Info("Frevod is waiting for paths or URIs in your terminal!")
	c.slog.Info("Use: {videoUri} {key} | eg: http://coolvideo.com home/coolvideo")

	var videoUri string
	var videoKeyStarter string
	for {

		_, err := fmt.Scanf("%v %v", &videoUri, &videoKeyStarter)
		if err != nil {
			c.slog.Error("Need URI and KEY | eg: /where/it/is where/to/put")
			continue
		}

		uriType := input_methods.CategorizeUri(videoUri)
		if uriType == "unsupported uri" {
			c.slog.Error("Unsupported URI!")
			c.slog.Info("Supported URIs starts with: http, https and file.")
			continue
		}

		job := core.NewJob()
		job.DownloadJob.Source = "cli"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = videoUri

		job.UploadJob.S3KeyStarter = videoKeyStarter

		c.downloadQueue <- job

		time.Sleep(time.Second) // avoid spam
	}
}
