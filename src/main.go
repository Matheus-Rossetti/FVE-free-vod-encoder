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

	filePool := make(chan *os.File, options.Encode.ConcurrentEncodings*2)
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)

	videoSlot := workspace.Prepare(options)
	for _, file := range videoSlot {
		filePool <- file
	}

	// START INPUT METHODS
	if options.Input.UseTerminal {
		go cli.Start(downloadQueue)
	}
	if options.Input.UseREST {
		go rest.Start(downloadQueue, ":8080")
	}

	// START DOWNLOADERS
	for index := range options.Encode.ConcurrentEncodings * 2 {
		go downloader.Start(index, downloadQueue, encodeQueue, filePool, options)
	}

	// START ENCODERS
	for index := range options.Encode.ConcurrentEncodings {
		go encoder.Start(index, encodeQueue, uploadQueue, filePool, options)
	}

	// START UPLOADER
	// Uploader manages concurrency inside
	// Do not start more than 1 uploader
	go uploader.Start(0, uploadQueue, options)

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
