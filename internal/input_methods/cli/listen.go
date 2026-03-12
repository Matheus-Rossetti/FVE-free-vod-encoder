package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (c *cli) ListenForInput() {
	c.slog.Info("Frevod is waiting for URIs in your terminal!")
	c.slog.Info("Use: {videoUri} {key}")
	c.slog.Info("Example: http://coolvideo.com home/coolvideo")

	var (
		videoUri        string
		videoKeyStarter string
	)

	for {
		_, err := fmt.Scanf("%v %v", &videoUri, &videoKeyStarter)
		if err != nil {
			c.slog.Error("Need URI and KEY!")
			c.slog.Error("Example: /where/it/is where/to/put")
			continue
		}

		uriType, err := core.CategorizeUri(videoUri)
		if errors.Is(err, core.ErrUnsupportedUri) {
			c.slog.Error("Unsupported URI!")
			c.slog.Info("Supported URIs starts with: http, https and file.")
			continue

		} else if err != nil {
			c.slog.Error("Failed to categorize URI", "error", err)
			continue
		}

		job := core.NewJob()
		job.DownloadJob.Source = "cli"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = videoUri

		job.UploadJob.KeyStarter = videoKeyStarter

		c.downloadQueue <- job

		time.Sleep(time.Second) // avoid spam
	}
}
