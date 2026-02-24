package core

import (
	"fmt"
)

func StartWorker(workerId int, jobQueue <-chan VideoJob, options *Options) {
	for job := range jobQueue {
		fmt.Printf("Worker %v received %v from %v", workerId, job.AbsoluteVideoPath, job.Source)

		video := NewVideo(job.AbsoluteVideoPath)
		outputDir := CreateOutputDir(video.Name, options)
		command := BuildFFmpegCommand(video, options)

		// FFmpeg runs as a low-priority process, it will use 100% CPU but won't freeze the system
		RunFFmpeg(command, outputDir)

		// TODO call interface to notify video finished
		Notification.FinishedEncoding()

		fmt.Printf("Job %v concluded!\n", video.Name)
	}
}
