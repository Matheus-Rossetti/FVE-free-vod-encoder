package downloader

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(
	ctx context.Context,
	options *core.Options,
	filePool chan *os.File,
	id int,
	downloadQueue <-chan *core.Job,
	encodeQueue chan<- *core.Job,
) {
JobLoop:
	for job := range downloadQueue {
		fmt.Printf("Download worker %v received %v from %v\n", id, job.DownloadJob.UriType, job.DownloadJob.Source)

		switch job.DownloadJob.UriType {
		case "url":
			file := <-filePool // file is returned to the pool by the encoder
			file.Truncate(0)
			file.Seek(0, 0) // 'Go' to the beginning of the file

			err := DownloadToFile(ctx, job.DownloadJob.VideoUri, file)
			if err != nil {
				fmt.Printf("\nError downloading %v to file", job.DownloadJob.VideoUri)
				filePool <- file
				continue JobLoop
			}
			job.EncodeJob.DownloadedFile = true
			job.EncodeJob.File = file

		case "path":
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

		encodeQueue <- job
	}

	// After queue closes
	fmt.Printf("\nDownloader %v shutting down...", id)
}
