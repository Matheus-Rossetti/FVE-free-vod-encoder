package downloader

import "github.com/Matheus-Rossetti/frevod/internal/core"

func Start(downloadQueue <-chan DownloadJob, options core.Options) {

	downloadSlots := make(chan DownloadJob, options.MaxStoredVideos)

	for downloadJob := range downloadQueue {
		downloadSlots <- downloadJob // when downloadSlots is full, this will wait until there's a free slot
		go DownloadAndStoreVideo(downloadJob.Uri)
	}
}
