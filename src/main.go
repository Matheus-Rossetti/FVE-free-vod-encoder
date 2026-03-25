package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/dispatcher"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/ingestor"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods/rest"
	"github.com/Matheus-Rossetti/frevod/internal/logger"
	"github.com/Matheus-Rossetti/frevod/internal/options"
	"github.com/Matheus-Rossetti/frevod/internal/output_methods/local"
	"github.com/Matheus-Rossetti/frevod/internal/output_methods/s3"
)

func main() {
	// WELLCOME TO THE FREVOD SOURCE CODE!

	// SLEEPS ARE SO LOGS COME OUT IN THE RIGHT ORDER

	frevod := core.StartFrevod(logger.Frevod())
	frevod.CheckForFFmpegBin()
	frevod.Greet()

	options := options.NewOptions(logger.Options())
	config := options.ParseOptions()

	// START FILE POOL (used by downloader and encoder)
	filePool := make(chan *os.File, config.Encode.ConcurrentEncodings*2)
	files := frevod.PrepareStorageFiles(config)
	frevod.FillFilePool(files, filePool)

	// START CONTEXT - LISTEN FOR SIGINT AND SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup

	// MAKE QUEUES
	downloadQueue := make(chan *core.Job, 999)
	encodeQueue := make(chan *core.Job, config.Encode.ConcurrentEncodings)
	uploadQueue := make(chan *core.Job, config.Encode.ConcurrentEncodings)
	wg.Go(func() {
		// The order in which the queues are closed is important
		// If we close the encoder when downloader is pushing a
		// job to it, the push will fail and leave orphan files.
		<-ctx.Done()
		frevod.Slog.Warn("Shutdown signal received!")

		time.Sleep(time.Second / 2)
		close(downloadQueue)

		time.Sleep(time.Second / 2)
		close(encodeQueue)

		time.Sleep(time.Second / 2)
		close(uploadQueue)
	})

	// START OUTPUT METHODS
	storageProviders := make(map[string]dispatcher.StorageProvider) // we pass storageProviders to uploader.Start
	if config.Upload.Local.Use {
		time.Sleep(time.Second / 2)
		log, slog := logger.Local()
		provider := local.Start(log, slog, config)
		storageProviders["local"] = provider
	}
	if config.Upload.S3.Use {
		time.Sleep(time.Second / 2)
		log, slog := logger.S3()
		provider := s3.Start(log, slog, config)
		storageProviders["s3"] = provider
	}

	// START UPLOAD
	time.Sleep(time.Second / 2)
	for index := range config.Encode.ConcurrentEncodings {
		wg.Go(func() {
			log, slog := logger.Ingestor(index)
			uploader := dispatcher.NewDispatcher(ctx, log, slog, config, index, uploadQueue, storageProviders)
			uploader.Start()
		})
	}

	// START ENCODERS
	time.Sleep(time.Second / 2)
	for index := range config.Encode.ConcurrentEncodings {
		wg.Go(func() {
			log, slog := logger.Encoder(index)
			encoder := encoder.NewEncoder(ctx, log, slog, config, filePool, index, encodeQueue, uploadQueue)
			encoder.Start()
		})
	}

	// START DOWNLOADERS
	time.Sleep(time.Second / 2)
	for index := range config.Encode.ConcurrentEncodings * 2 {
		wg.Go(func() {
			log, slog := logger.Dispatcher(index)
			downloader := ingestor.NewIngestor(ctx, log, slog, config, filePool, index, downloadQueue, encodeQueue)
			downloader.Start()
		})

	}

	// START INPUT METHODS
	if config.Input.Cli {
		time.Sleep(time.Second / 2)
		log, slog := logger.Cli()
		wg.Go(func() {
			cli := cli.NewCli(ctx, log, slog, downloadQueue)
			cli.Start()
		})
	}
	if config.Input.REST {
		time.Sleep(time.Second / 2)
		log, slog := logger.Rest()
		wg.Go(func() {
			rest := rest.NewRest(ctx, log, slog, ":8080", downloadQueue)
			rest.Start()
		})
	}

	// Shutdown
	wg.Wait()
	close(filePool)
	frevod.CloseAndDeleteStorageFiles(filePool)
	os.RemoveAll("storage_files")
	stop() // ctx

	frevod.Slog.Info("Everything's clean, bye :)")
}
