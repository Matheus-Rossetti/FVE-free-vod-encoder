package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/cli"
	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func main() {

	core.CheckForFFmpegBin()
	options := core.ParseOptions()
	core.Greet()

	jobQueue := make(chan core.VideoJob, 999)

	if options.UseTerminal {
		go cli.Start(jobQueue)
	}

	for workerId := range options.ConcurrentEncodings {
		go core.StartWorker(workerId, jobQueue, options)
	}

	for {
		var quit string
		fmt.Scanln(quit)
		if quit == "q" {
			break
		}
	}

}
