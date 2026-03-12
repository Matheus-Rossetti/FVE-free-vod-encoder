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
		d.slog.Info(
			fmt.Sprintf("Received %v from %v", job.DownloadJob.UriType.String(), job.DownloadJob.Source),
			"id", d.id,
		)

		switch job.DownloadJob.UriType {
		case core.Url:
			file := <-d.filePool // file is returned to the pool by the encoder or by an error

			err := d.prepareFileForDownload(file)
			if err != nil {
				d.slog.Error(ErrPreparingFileforDownload.Error(), "err", err)
				continue JobLoop
			}

			err = d.downloadToFile(job.DownloadJob.VideoUri, file)
			if err != nil {
				d.slog.Error(ErrDownloading.Error(), "url", job.DownloadJob.VideoUri, "err", err)
				d.filePool <- file
				continue JobLoop

			}
			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case core.Path:
			localFile, err := d.validateFileFromPath(job.DownloadJob.VideoUri)
			if err != nil {
				d.slog.Error(ErrValidatingFileFromPath.Error(), "err", err)
				continue JobLoop
			}

			job.EncodeJob.DownloadedFile = false
			job.EncodeJob.File = localFile
		}

		d.encodeQueue <- job
	}

	// After queue closes
	d.slog.Warn("Shutting down...", "id", d.id)
}
