package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/downloader"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/rest"
	"github.com/Matheus-Rossetti/frevod/internal/logger"
	"github.com/Matheus-Rossetti/frevod/internal/options"
	"github.com/Matheus-Rossetti/frevod/internal/output_methods/local"
	"github.com/Matheus-Rossetti/frevod/internal/output_methods/s3"
	"github.com/Matheus-Rossetti/frevod/internal/uploader"
)

func main() {

	frevod := StartFrevod(logger.Frevod())
	frevod.CheckForFFmpegBin()
	frevod.Greet()

	options := options.NewOptions(logger.Options())
	config := options.ParseOptions()

	// MAKE QUEUES
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, config.Encode.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, config.Encode.ConcurrentEncodings)

	// START CONTEXT - LISTEN FOR SIGINT AND SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		frevod.slog.Info("Shutdown signal received")
		close(downloadQueue)
		close(encodeQueue)
		close(uploadQueue)
	}()

	var wg sync.WaitGroup

	// START INPUT METHODS
	if config.Input.UseTerminal {
		log, slog := logger.Cli()
		wg.Go(func() {
			cli := cli.NewCli(ctx, log, slog, downloadQueue)
			cli.Start()
		})
	}
	if config.Input.UseREST {
		log, slog := logger.Rest()
		wg.Go(func() {
			rest := rest.NewRest(ctx, log, slog, "8080", downloadQueue)
			rest.Start()
		})
	}

	// START OUTPUT METHODS
	storageProviders := make(map[string]uploader.StorageProvider) // we pass storageProviders to uploader.Start
	if config.Upload.Local.Use {
		provider := local.Start(config)
		storageProviders["local"] = provider
	}
	if config.Upload.S3.Use {
		provider := s3.Start(config)
		storageProviders["s3"] = provider
	}

	// START FILE POOL (used by downloader and encoder)
	filePool := make(chan *os.File, config.Encode.ConcurrentEncodings*2)
	files := frevod.PrepareStorageFiles(config)
	frevod.FillFilePool(files, filePool)

	// START DOWNLOADERS
	for index := range config.Encode.ConcurrentEncodings * 2 {
		wg.Add(1)
		go func() {
			// this doesn't work, maybe because when you pass a chan to a struct
			// it becomes another channel, and not the same
			// the downloadQueue doesn't close if we pass it into the struct
			log, slog := logger.Downloader()
			downloader := downloader.NewDownloader(log, slog, ctx, config, filePool, index, downloadQueue, encodeQueue)
			downloader.Start()
			wg.Done()
		}()
	}

	// START ENCODERS
	for index := range config.Encode.ConcurrentEncodings {
		wg.Add(1)
		go func() {
			encoder.Start(ctx, config, filePool, index, encodeQueue, uploadQueue)
			wg.Done()
		}()
	}

	// START UPLOAD
	for index := range config.Encode.ConcurrentEncodings {
		wg.Add(1)
		go func() {
			uploader.Start(ctx, config, index, uploadQueue, storageProviders)
			wg.Done()
		}()
	}

	// Shutdown
	wg.Wait()
	close(filePool)
	frevod.CloseAndDeleteStorageFiles(filePool)
	os.RemoveAll("storage_files")
	stop()

	fmt.Printf("\n\nEverything's clean, bye :)\n")
}
