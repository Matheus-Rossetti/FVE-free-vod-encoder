package ingestor

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type ingestor struct {
	ctx           context.Context
	log           *log.Logger
	slog          *slog.Logger
	options       *core.Options
	filePool      chan *os.File
	id            int
	downloadQueue <-chan *core.Job
	encodeQueue   chan<- *core.Job
}

func NewIngestor(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	options *core.Options,
	filePool chan *os.File,
	id int,
	downloadQueue <-chan *core.Job,
	encodeQueue chan<- *core.Job) *ingestor {

	return &ingestor{
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
