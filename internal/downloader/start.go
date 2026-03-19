package downloader

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

func (d *downloader) Start() {
JobLoop:
	for job := range d.downloadQueue {
		d.slog.Info(fmt.Sprintf("Received a job from %v", job.DownloadJob.Source),
			"id", d.id)

		switch job.DownloadJob.UriType {
		case core.Url:
			file := <-d.filePool // file is returned to the pool by the encoder or by an error
			d.slog.Info(fmt.Sprintf("Downloading from %v into %v", job.DownloadJob.VideoUri, file.Name()),
				"id", d.id)

			err := d.prepareFileForDownload(file)
			if err != nil {
				d.slog.Error(ErrPreparingFileforDownload.Error(), "err", err, "id", d.id)
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
