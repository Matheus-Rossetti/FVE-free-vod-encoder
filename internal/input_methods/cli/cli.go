package cli

import (
	"context"
	"log"
	"log/slog"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type cli struct {
	ctx           context.Context
	log           *log.Logger
	slog          *slog.Logger
	downloadQueue chan<- *core.Job
	err           error
}

func NewCli(ctx context.Context, log *log.Logger, slog *slog.Logger, downloadQueue chan<- *core.Job) *cli {
	return &cli{
		ctx:           ctx,
		log:           log,
		slog:          slog,
		downloadQueue: downloadQueue,
		err:           nil,
	}
}
