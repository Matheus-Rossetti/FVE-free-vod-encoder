package main

import (
	"fmt"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/cli"
	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func main() {

	core.CheckForFFmpegBin()
	core.Greet()

	options := core.ParseOptions()

	jobQueue := make(chan core.VideoJob, 999)

	if options.UseTerminal {
		go cli.Start(jobQueue)
	}

	for job := range jobQueue {
		go func() {
			start := time.Now()
			video := core.NewVideo(job.AbsoluteVideoPath)
			outputDir := core.CreateOutputDir(video.Name)
			command := core.BuildFFmpegCommand(video, options)
			core.RunFFmpeg(command, outputDir)
			fmt.Printf("Job %v concluded!\n", video.Name)
			fmt.Printf("\n\n Encoding took %v", time.Since(start))
		}()
	}
}
