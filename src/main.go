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
	"github.com/Matheus-Rossetti/frevod/internal/output_methods/s3"
	"github.com/Matheus-Rossetti/frevod/internal/uploader"
)

func main() {

	core.CheckForFFmpegBin()
	options := options.ParseOptions()
	core.Greet()

	// MAKE QUEUES
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, options.Encode.ConcurrentEncodings)

	// START CONTEXT - LISTEN FOR SIGINT AND SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		close(downloadQueue)
		close(encodeQueue)
		close(uploadQueue)
	}()

	var wg sync.WaitGroup

	// START INPUT METHODS
	if options.Input.UseTerminal {
		go cli.Start(ctx, downloadQueue)
	}
	if options.Input.UseREST {
		go rest.Start(ctx, downloadQueue, ":8080")
	}

	// START OUTPUT METHODS
	storageProviders := make(map[string]uploader.StorageProvider)
	if options.Upload.S3.Use {
		provider := s3.Start(options)
		storageProviders["s3"] = provider
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

	// START UPLOAD
	for range options.Encode.ConcurrentEncodings {
		go uploader.Start(ctx, options, uploadQueue, storageProviders)
	}

	// Shutdown
	wg.Wait()
	close(filePool)
	core.CloseAndDeleteStorageFiles(filePool)
}
