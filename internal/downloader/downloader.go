package downloader

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type downloader struct {
	ctx           context.Context
	log           *log.Logger
	slog          *slog.Logger
	options       *core.Options
	filePool      chan *os.File
	id            int
	downloadQueue <-chan *core.Job
	encodeQueue   chan<- *core.Job
}

func NewDownloader(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	options *core.Options,
	filePool chan *os.File,
	id int,
	downloadQueue <-chan *core.Job,
	encodeQueue chan<- *core.Job) *downloader {

	return &downloader{
		ctx:           ctx,
		log:           log,
		slog:          slog,
		options:       options,
		filePool:      filePool,
		id:            id,
		downloadQueue: downloadQueue,
		encodeQueue:   encodeQueue,
	}
}
