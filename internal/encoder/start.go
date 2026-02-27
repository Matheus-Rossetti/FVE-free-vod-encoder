package encoder

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(workerId int, encodeQueue <-chan *core.Job, uploadQueue chan<- *core.Job, filePool chan<- *os.File, options *core.Options) {
	for job := range encodeQueue {
		fmt.Printf("Encode worker %v received %v\n", workerId, job.EncodeJob.AbsoluteVideoPath)

		video := NewVideo(job.EncodeJob.AbsoluteVideoPath)
		cmd := BuildFFmpegCommand(video, options)

		outputDir := createOutputDir()
		// FFmpeg runs as a low-priority process, it will use 100% CPU but won't freeze the system
		cmd.Dir = outputDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal("error running the command\n", err, "for:", string(output))
		}

		filePool <- job.EncodeJob.File // return the file to the pool

		job.UploadJob.DirToUploadFrom = outputDir
		uploadQueue <- job

		fmt.Printf("Job %v concluded!\n", video.Name)
	}
}
