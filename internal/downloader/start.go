package downloader

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(downloaderId int, downloadQueue <-chan core.Job, encodeQueue chan<- core.Job, options *core.Options) {
	for job := range downloadQueue {
		fmt.Printf("Download worker %v received %v from %v\n", downloaderId, job.VideoUri, job.Source)

		videoPath := DownloadToFile(job.VideoUri, job.File)
		job.SetAbsolutePath(videoPath)

		encodeQueue <- job
	}
}
