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

func (d *ingestor) Start() {
JobLoop:
	for job := range d.downloadQueue {
		d.slog.Info("Received a job!", "source", job.DownloadJob.Source)

		switch job.DownloadJob.UriType {
		case core.Url:
			file := <-d.filePool // file is returned to the pool by the encoder or by an error
			d.slog.Info("Downloading...", "from", job.DownloadJob.VideoUri)

			err := d.prepareFileForDownload(file)
			if err != nil {
				d.slog.Error(ErrPreparingFileforDownload.Error(), "err", err)
				d.filePool <- file
				continue JobLoop
			}

			err = d.downloadToFile(job.DownloadJob.VideoUri, file)
			if err != nil {
				d.slog.Error(ErrDownloading.Error(), "url", job.DownloadJob.VideoUri, "err", err, "id", d.id)
				d.filePool <- file
				continue JobLoop
			}

			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case core.Path:
			d.slog.Info(fmt.Sprintf("Validating %v", job.DownloadJob.VideoUri),
				"id", d.id)

			localFile, err := d.validateFileFromPath(job.DownloadJob.VideoUri)
			if err != nil {
				d.slog.Error(ErrValidatingFileFromPath.Error(), "err", err, "id", d.id)
				continue JobLoop
			}

			job.EncodeJob.DownloadedFile = false
			job.EncodeJob.File = localFile
		}

		d.slog.Info("Finished!", "id", d.id)
		d.encodeQueue <- job
	}

	// After queue closes
	d.slog.Info("Shutting down...", "id", d.id)
}
