package encoder

import (
	"fmt"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(workerId int, encodeQueue <-chan core.Job, options *core.Options, notifier INotifier) {
	for job := range encodeQueue {
		fmt.Printf("Encode worker %v received %v", workerId, job.AbsoluteVideoPath)

		video := NewVideo(job.AbsoluteVideoPath)
		// outputDir := CreateOutputDir(video.Name, options)
		cmd := BuildFFmpegCommand(video, options)

		// FFmpeg runs as a low-priority process, it will use 100% CPU but won't freeze the system
		cmd.Dir = job.OutputDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal("error running the command\n", err, "for:", string(output))
		}

		// TODO call interface to notify video finished
		notifier.FinishedEncoding()

		fmt.Printf("Job %v concluded!\n", video.Name)
	}
}
