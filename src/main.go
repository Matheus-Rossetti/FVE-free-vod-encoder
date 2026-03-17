package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	// HELLO!
	// THIS IS THE MAIN FUNCTION OF FREVOD
	// IT ORCHESTRATES EVERYTHING

	// SLEEPS ARE SO LOGS COME OUT IN THE RIGHT ORDER

	frevod := startFrevod(logger.Frevod())
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
		// The order in which the queues are closed is important
		// If we close encoder is closes when downloader finishes a download job
		// the downloader will try to push the job to the encoder, fail and leave orphan files.
		<-ctx.Done()
		frevod.slog.Warn("Shutdown signal received!")

		time.Sleep(time.Second / 2)
		close(downloadQueue)

		time.Sleep(time.Second / 2)
		close(encodeQueue)

		time.Sleep(time.Second / 2)
		close(uploadQueue)
	}()

	var wg sync.WaitGroup

	// START INPUT METHODS
	if config.Input.UseTerminal {
		time.Sleep(time.Second / 2)
		log, slog := logger.Cli()
		wg.Go(func() {
			cli := cli.NewCli(ctx, log, slog, downloadQueue)
			cli.Start()
		})
	}
	if config.Input.UseREST {
		time.Sleep(time.Second / 2)
		log, slog := logger.Rest()
		wg.Go(func() {
			rest := rest.NewRest(ctx, log, slog, "8080", downloadQueue)
			rest.Start()
		})
	}

	// START OUTPUT METHODS
	storageProviders := make(map[string]uploader.StorageProvider) // we pass storageProviders to uploader.Start
	if config.Upload.Local.Use {
		time.Sleep(time.Second / 2)
		provider := local.Start(config)
		storageProviders["local"] = provider
	}
	if config.Upload.S3.Use {
		time.Sleep(time.Second / 2)
		provider := s3.Start(config)
		storageProviders["s3"] = provider
	}

	// START FILE POOL (used by downloader and encoder)
	filePool := make(chan *os.File, config.Encode.ConcurrentEncodings*2)
	files := frevod.PrepareStorageFiles(config)
	frevod.FillFilePool(files, filePool)

	// START DOWNLOADERS
	time.Sleep(time.Second / 2)
	for index := range config.Encode.ConcurrentEncodings * 2 {
		wg.Go(func() {
			log, slog := logger.Downloader()
			downloader := downloader.NewDownloader(ctx, log, slog, config, filePool, index, downloadQueue, encodeQueue)
			downloader.Start()
		})
	}

	// START ENCODERS
	time.Sleep(time.Second / 2)
	for index := range config.Encode.ConcurrentEncodings {
		wg.Go(func() {
			log, slog := logger.Encoder()
			encoder := encoder.NewEncoder(ctx, log, slog, config, filePool, index, encodeQueue, uploadQueue)
			encoder.Start()
		})
	}

	// START UPLOAD
	time.Sleep(time.Second / 2)
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
