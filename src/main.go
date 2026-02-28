package main

import (
	"fmt"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/rest"
	"github.com/Matheus-Rossetti/frevod/internal/options"
	"github.com/Matheus-Rossetti/frevod/internal/uploader"
	"github.com/Matheus-Rossetti/frevod/internal/workspace"
)

func main() {

	core.CheckForFFmpegBin()
	options := options.ParseOptions()
	core.Greet()

	filePool := make(chan *os.File, options.ConcurrentEncodings*2)
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, options.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, options.ConcurrentEncodings)

	videoSlot := workspace.Prepare(options)
	for _, file := range videoSlot {
		filePool <- file
	}

	// START DOWNLOADERS
	for index := range options.ConcurrentEncodings * 2 {
		go downloader.Start(index, downloadQueue, encodeQueue, filePool, options)
	}

	// START ENCODERS
	for index := range options.ConcurrentEncodings {
		go encoder.Start(index, encodeQueue, uploadQueue, filePool, options)
	}

	// START UPLOADERS
	go uploader.Start(0, uploadQueue) // just one for testing

	// START INPUT METHODS
	if options.UseTerminal {
		go cli.Start(downloadQueue)
	}
	if options.UseREST {
		go rest.Start(downloadQueue, ":8080")
	}

	for {
		var quit string
		fmt.Scanln(quit)
		if quit == "q" {
			for _, file := range videoSlot {
				file.Close()
			}
			break
		}
	}
}
