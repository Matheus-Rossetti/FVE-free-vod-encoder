package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/cli"
	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
)

func main() {

	core.CheckForFFmpegBin()
	options := core.ParseOptions()
	core.Greet()

	jobQueue := make(chan core.VideoJob, 999)
	downloadQueue := make(chan downloader.DownloadJob, 999)

	if options.UseTerminal {
		go cli.Start(downloadQueue)
	}
	// if options.UseREST {
	// 	go rest.Start(downloadQueue, ":8080")
	// }

	downloader.Start(downloadQueue, *options)

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
