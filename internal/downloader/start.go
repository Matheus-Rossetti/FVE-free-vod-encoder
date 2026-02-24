package downloader

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(downloaderId int, downloadJobQueue <-chan DownloadJob, videoJobQueue chan<- core.VideoJob, options *core.Options) {
	for downloadJob := range downloadJobQueue {
		fmt.Printf("Download worker %v received %v from %v\n", downloaderId, downloadJob.Uri, downloadJob.Source)
		videoPath := DownloadAndStoreVideo(downloadJob.Uri)
		videoJobQueue <- core.VideoJob{
			AbsoluteVideoPath: videoPath,
		}
	}
}
