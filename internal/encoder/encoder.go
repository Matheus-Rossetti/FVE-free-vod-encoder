package encoder

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type encoder struct {
	ctx           context.Context
	log           *log.Logger
	slog          *slog.Logger
	options       *core.Options
	filePool      chan<- *os.File
	encodeQueue   <-chan *core.Job
	dispatchQueue chan<- *core.Job
}

func NewEncoder(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	options *core.Options,
	filePool chan<- *os.File,
	encodeQueue core.EncodeQueue,
	dispatchQueue core.DispatchQueue) *encoder {
	return &encoder{
		ctx:           ctx,
		log:           log,
		slog:          slog,
		options:       options,
		filePool:      filePool,
		encodeQueue:   encodeQueue,
		dispatchQueue: dispatchQueue,
	}
}
