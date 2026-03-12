package downloader

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type downloader struct {
	log           *log.Logger
	slog          *slog.Logger
	ctx           context.Context
	options       *core.Options
	filePool      chan *os.File
	id            int
	downloadQueue <-chan *core.Job
	encodeQueue   chan<- *core.Job
}

func NewDownloader(
	log *log.Logger,
	slog *slog.Logger,
	ctx context.Context,
	options *core.Options,
	filePool chan *os.File,
	id int,
	downloadQueue <-chan *core.Job,
	encodeQueue chan<- *core.Job) *downloader {

	return &downloader{
		log:           log,
		slog:          slog,
		ctx:           ctx,
		options:       options,
		filePool:      filePool,
		id:            id,
		downloadQueue: downloadQueue,
		encodeQueue:   encodeQueue,
	}
}
