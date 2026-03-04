package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/rest"
	"github.com/Matheus-Rossetti/frevod/internal/options"
	"github.com/Matheus-Rossetti/frevod/internal/uploader"
)

func main() {

	core.CheckForFFmpegBin()
	options := options.ParseOptions()
	core.Greet()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	// MAKE QUEUES
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)

	// START INPUT METHODS
	if options.Input.UseTerminal {
		go cli.Start(downloadQueue)
	}
	if options.Input.UseREST {
		go rest.Start(downloadQueue, ":8080")
	}

	// START FILE POOL (used by downloader and encoder)
	filePool := make(chan *os.File, options.Encode.ConcurrentEncodings*2)
	files := core.PrepareStorageFiles(options)
	core.FillFilePool(files, filePool)

	// START DOWNLOADERS
	for index := range options.Encode.ConcurrentEncodings * 2 {
		wg.Add(1)
		go func() {
			downloader.Start(ctx, options, filePool, index, downloadQueue, encodeQueue)
			wg.Done()
		}()
	}

	// START ENCODERS
	for index := range options.Encode.ConcurrentEncodings {
		wg.Add(1)
		go func() {
			encoder.Start(ctx, options, filePool, index, encodeQueue, uploadQueue)
			wg.Done()
		}()
	}

	// START UPLOADER
	// Uploader manages concurrency inside
	// Do not start more than 1 uploader
	go uploader.Start(0, uploadQueue, options)

	// Shutdown
	wg.Wait()
	close(filePool)
	core.CloseAndDeleteStorageFiles(filePool)
}
