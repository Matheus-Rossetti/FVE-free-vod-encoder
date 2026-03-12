package downloader

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (d *downloader) Start() {
JobLoop:
	for job := range d.downloadQueue {
		d.slog.Info(
			fmt.Sprintf("Received %v from %v", job.DownloadJob.UriType, job.DownloadJob.Source),
			"id", d.id,
		)

		switch job.DownloadJob.UriType {
		case core.Url:
			file := <-d.filePool // file is returned to the pool by the encoder
			file.Truncate(0)
			file.Seek(0, 0) // 'Go' to the beginning of the file

			err := d.DownloadToFile(d.ctx, job.DownloadJob.VideoUri, file)
			if err != nil {
				fmt.Printf("\nError downloading %v to file", job.DownloadJob.VideoUri)
				d.filePool <- file
				continue JobLoop
			}
			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case core.Path:
			strings.TrimPrefix(job.DownloadJob.VideoUri, "file://")
			localFile, err := os.Open(job.DownloadJob.VideoUri)
			defer localFile.Close()
			if err != nil {
				fmt.Printf("\nError opening %v", job.DownloadJob.VideoUri)
				continue JobLoop
			}

			// check if it's video
			buffer := make([]byte, 512)
			bytesRead, err := localFile.Read(buffer)
			if err != nil {
				fmt.Printf("\nError reading %v", job.DownloadJob.VideoUri)
				continue JobLoop
			}
			MIMEtype := http.DetectContentType(buffer[:bytesRead])
			fileType := strings.Split(MIMEtype, "/")
			if fileType[0] != "video" {
				fmt.Printf("\nError: File isn't a video, it's %v", fileType[0])
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
