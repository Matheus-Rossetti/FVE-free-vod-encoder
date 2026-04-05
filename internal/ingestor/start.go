package ingestor

import (
	"errors"
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

var (
	ErrPreparingFileforDownload = errors.New("failed to prepare file")
	ErrDownloading              = errors.New("failed to download file")
	ErrValidatingFileFromPath   = errors.New("failed local file validation")
)

func (i *ingestor) Start() {
JobLoop:
	for job := range i.downloadQueue {
		i.slog.Info("Received a job!", "source", job.DownloadJob.Source)

		switch job.DownloadJob.UriType {
		case core.Url:
			file := <-i.filePool // file is returned to the pool by the encoder or by an error

			err := i.ingestFromUrl(file, job.DownloadJob.VideoUri)
			if err != nil {
				i.filePool <- file
				i.handleUrlError(err)
				continue
			}

			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case core.Path:
			i.slog.Info(fmt.Sprintf("Validating %v", job.DownloadJob.VideoUri))

			localFile, err := i.validateFileFromPath(job.DownloadJob.VideoUri)
			if err != nil {
				i.slog.Error(ErrValidatingFileFromPath.Error(), "err", err)
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
