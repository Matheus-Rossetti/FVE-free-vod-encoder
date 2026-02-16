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
			outputDir := core.CreateOutputDir(job.Id)
			command := core.BuildFFmpegCommand(job.VideoPath, "video1")
			core.RunFFmpeg(command, outputDir)
			fmt.Printf("Job %v concluded! ", job.Id)
		}()
	}
}
