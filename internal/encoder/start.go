package encoder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(
	ctx context.Context,
	options *core.Options,
	filePool chan<- *os.File,
	id int,
	encodeQueue <-chan *core.Job,
	uploadQueue chan<- *core.Job,
) {
JobLoop:
	for job := range encodeQueue {
		defer func() {
			if job.EncodeJob.DownloadedFile {
				filePool <- job.EncodeJob.File
			}
		}() // return the file to the pool

		fmt.Printf("Encode worker %v received %v\n", id, job.EncodeJob.AbsoluteVideoPath)
		outputDir := createOutputDir()

		absoluteVideoPath, _ := filepath.Abs(job.EncodeJob.File.Name())

		video := NewVideo(absoluteVideoPath)
		cmd := BuildFFmpegCommand(ctx, options, video)

		cmd.Dir = outputDir
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error while running the command: %v\nDeleting: %v", err, outputDir)
			os.RemoveAll(outputDir)
			continue JobLoop
		}

		job.UploadJob.FromDir = outputDir
		uploadQueue <- job

		fmt.Printf("Job %v concluded!\n", video.Name)
	}

	// After queue closes
	os.RemoveAll("output")
	fmt.Printf("\nEncode %v shuting down...", id)
}
