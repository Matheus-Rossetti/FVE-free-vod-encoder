package downloader

import (
	"context"
	"fmt"
	"os"

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
	// We need to repeat the case <-ctx.Done(),
	// otherwise, if the second select is waiting
	// for something from the downloadQueue only,
	// it will stay stuck there and never
	// shut down, even if done is closed
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Downloader %v shutting down...\n", id)
			return

		default:
		}

		select {
		case <-ctx.Done():
			fmt.Printf("Downloader %v shutting down...\n", id)
			return

		case job := <-downloadQueue:
			file := <-filePool // file is returned to the pool by the encoder after it's done encoding said file

			fmt.Printf("Download worker %v received %v from %v\n", id, job.DownloadJob.VideoUri, job.DownloadJob.Source)

			file.Truncate(0) // Clean the file without deleting it
			file.Seek(0, 0)  // 'Go' to the beginning of the file
			err := DownloadToFile(ctx, job.DownloadJob.VideoUri, file)
			if err != nil {
				filePool <- file
				break
			}

			job.EncodeJob.File = file
			encodeQueue <- job
		}
	}
}
