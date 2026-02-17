package main

import (
	"fmt"

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
			video := core.NewVideo(job.AbsoluteVideoPath)
			outputDir := core.CreateOutputDir(video)
			command := core.BuildFFmpegCommand(video)
			core.RunFFmpeg(command, outputDir)
			fmt.Printf("Job %v concluded!\n", video.Name)
		}()
	}
}
