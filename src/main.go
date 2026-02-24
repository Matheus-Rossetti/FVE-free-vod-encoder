package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/cli"
	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
	"github.com/Matheus-Rossetti/frevod/internal/notifier"
)

func main() {

	core.CheckForFFmpegBin()
	options := core.ParseOptions()
	core.Greet()

	downloadJobQueue := make(chan downloader.DownloadJob, 999)
	videoJobQueue := make(chan core.VideoJob, options.ConcurrentEncodings)

	notifier := notifier.New()

	if options.UseTerminal {
		go cli.Start(downloadJobQueue)
	}
	// if options.UseREST {
	// 	go rest.Start(downloadQueue, ":8080")
	// }

	for downloaderId := range options.ConcurrentEncodings {
		go downloader.Start(downloaderId, downloadJobQueue, videoJobQueue, options)
	}

	for processorId := range options.ConcurrentEncodings {
		go core.Start(processorId, videoJobQueue, options, notifier)
	}

	for {
		var quit string
		fmt.Scanln(quit)
		if quit == "q" {
			break
		}
	}
}
