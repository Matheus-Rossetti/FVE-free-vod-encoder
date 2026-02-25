package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/notifier"
	"github.com/Matheus-Rossetti/frevod/internal/options"
)

func main() {

	core.CheckForFFmpegBin()
	options := options.ParseOptions()
	core.Greet()

	downloadQueue := make(chan core.Job, 999)
	encodeQueue := make(chan core.Job, options.ConcurrentEncodings)

	notifier := notifier.New()

	for downloaderId := range options.ConcurrentEncodings {
		go downloader.Start(downloaderId, downloadQueue, encodeQueue, options)
	}

	for processorId := range options.ConcurrentEncodings {
		go encoder.Start(processorId, encodeQueue, options, notifier)
	}

	if options.UseTerminal {
		go cli.Start(downloadQueue)
	}
	// if options.UseREST {
	// 	go rest.Start(downloadQueue, ":8080")
	// }

	for {
		var quit string
		fmt.Scanln(quit)
		if quit == "q" {
			break
		}
	}
}
