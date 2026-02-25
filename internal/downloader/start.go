package downloader

import (
	"fmt"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(id int, downloadQueue <-chan core.DownloadJob, encodeQueue chan<- core.EncodeJob, filePool <-chan *os.File, options *core.Options) {
	for downloadJob := range downloadQueue {
		file := <-filePool

		fmt.Printf("Download worker %v received %v from %v\n", id, downloadJob.VideoUri, downloadJob.Source)

		file.Truncate(0) // Cleans the file without deleting it
		file.Seek(0, 0)  // 'Go' to the beginning of the file
		DownloadToFile(downloadJob.VideoUri, file)

		encodeJob := core.NewEncodeJob(file)
		encodeQueue <- encodeJob
	}
}
