package cli

import (
	"context"
	"log"
	"log/slog"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type cli struct {
	ctx         context.Context
	log         *log.Logger
	slog        *slog.Logger
	ingestQueue core.IngestQueue
}

func NewCli(ctx context.Context, log *log.Logger, slog *slog.Logger, downloadQueue core.IngestQueue) *cli {
	return &cli{
		ctx:         ctx,
		log:         log,
		slog:        slog,
		ingestQueue: downloadQueue,
	}
}
