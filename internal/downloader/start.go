package downloader

import (
	"fmt"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(id int, downloadQueue <-chan *core.Job, encodeQueue chan<- *core.Job, filePool <-chan *os.File, options *core.Options) {
	for job := range downloadQueue {
		file := <-filePool

		fmt.Printf("Download worker %v received %v from %v\n", id, job.DownloadJob.VideoUri, job.DownloadJob.Source)

		file.Truncate(0) // Cleans the file without deleting it
		file.Seek(0, 0)  // 'Go' to the beginning of the file
		DownloadToFile(job.DownloadJob.VideoUri, file)

		job.EncodeJob.File = file
		encodeQueue <- job
	}
}
