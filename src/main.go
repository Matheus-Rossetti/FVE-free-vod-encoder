package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Matheus-Rossetti/frevod/internal/config"
	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/dispatch_methods/local"
	"github.com/Matheus-Rossetti/frevod/internal/dispatch_methods/s3"
	"github.com/Matheus-Rossetti/frevod/internal/dispatcher"
	"github.com/Matheus-Rossetti/frevod/internal/encoder"
	"github.com/Matheus-Rossetti/frevod/internal/ingest_methods/cli"
	"github.com/Matheus-Rossetti/frevod/internal/ingest_methods/rest"
	"github.com/Matheus-Rossetti/frevod/internal/ingestor"
	"github.com/Matheus-Rossetti/frevod/internal/logger"
)

func main() {
	// WELCOME TO THE FREVOD SOURCE CODE!

	frevod := core.StartFrevod(logger.Frevod())
	frevod.CheckForFFmpegBin()
	frevod.Greet()

	options := core.NewOptions()
	err := config.Load(options, "")
	if err != nil {
		frevod.Log.Fatal(err.Error())
	}

	// START FILE POOL (used by downloader and encoder)
	filePool := make(chan *os.File, options.Encode.ConcurrentEncodings*2)
	files := frevod.PrepareStorageFiles(options)
	frevod.FillFilePool(files, filePool)

	// START CONTEXT - LISTEN FOR SIGINT AND SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup

	ingestQueue, encodeQueue, dispatchQueue := core.StartQueues(options)

	wg.Go(func() {
		// The order in which the queues are closed is important
		// If we close the encoder when downloader is pushing a
		// job to it, the push will fail and leave orphan files.
		<-ctx.Done()
		frevod.Slog.Warn("Shutdown signal received!")
		close(ingestQueue)
		close(encodeQueue)
		close(dispatchQueue)
	})

	// START OUTPUT METHODS
	var storageProviders []dispatcher.StorageProvider // we pass storageProviders to dispatcher
	if options.Dispatch.S3.Enabled {
		log, slog := logger.S3()
		provider := s3.Start(log, slog, options)
		storageProviders = append(storageProviders, provider)
	}

	// local should be added last, since it renames the output dir.
	if options.Dispatch.Local.Enabled {
		log, slog := logger.Local()
		provider := local.Start(log, slog, options)
		storageProviders = append(storageProviders, provider)
	}

	//-------------------------------------------
	for index := range options.Encode.ConcurrentEncodings {
		wg.Go(func() {
			log, slog := logger.Dispatcher(index)
			uploader := dispatcher.NewDispatcher(ctx, log, slog, options, dispatchQueue, storageProviders)
			uploader.Start()
		})
	}

	for index := range options.Encode.ConcurrentEncodings {
		wg.Go(func() {
			log, slog := logger.Encoder(index)
			encoder := encoder.NewEncoder(ctx, log, slog, options, filePool, encodeQueue, dispatchQueue)
			encoder.Start()
		})
	}

	for index := range options.Encode.ConcurrentEncodings * 2 {
		wg.Go(func() {
			log, slog := logger.Ingestor(index)
			downloader := ingestor.NewIngestor(ctx, log, slog, options, filePool, index, ingestQueue, encodeQueue)
			downloader.Start()
		})

	}

	// START INPUT METHODS
	if options.Ingest.Cli.Enabled {
		log, slog := logger.Cli()
		wg.Go(func() {
			cli := cli.NewCli(ctx, log, slog)
			cli.Start()
		})
	}
	if options.Ingest.REST.Enabled {
		log, slog := logger.Rest()
		wg.Go(func() {
			rest := rest.NewRest(ctx, log, slog, options.Ingest.REST.Port)
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
