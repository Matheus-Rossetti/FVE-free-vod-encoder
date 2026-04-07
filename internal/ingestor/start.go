package ingestor

import (
	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (i *ingestor) Start() {
JobLoop:
	for job := range i.downloadQueue {
		i.slog.Info("Received a job!", "source", job.DownloadJob.Source)

		switch job.DownloadJob.UriType {
		case core.Url:
			i.slog.Info("Downloading...", "from", job.DownloadJob.VideoUri)

			file := <-i.filePool // file is returned to the pool by the encoder or by an error
			err := i.ingestFromUrl(file, job.DownloadJob.VideoUri)
			if err != nil {
				i.filePool <- file
				i.slog.Error("skipping job", "why", err)
				continue JobLoop
			}

			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case core.Path:
			i.slog.Info("Validating...", "file", job.DownloadJob.VideoUri)

			localFile, err := i.ingestFromLocal(job.DownloadJob.VideoUri)
			if err != nil {
				i.slog.Error("skipping job", "why", err)
				continue JobLoop
			}

			job.EncodeJob.DownloadedFile = false
			job.EncodeJob.File = localFile
		}

		i.slog.Info("Finished!")
		i.encodeQueue <- job
	}

	// After queue closes
	i.slog.Info("Shutting down...")
}
