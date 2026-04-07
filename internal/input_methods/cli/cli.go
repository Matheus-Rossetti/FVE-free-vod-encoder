package cli

import (
	"context"
	"log"
	"log/slog"
)

type cli struct {
	ctx  context.Context
	log  *log.Logger
	slog *slog.Logger
}

func NewCli(ctx context.Context, log *log.Logger, slog *slog.Logger) *cli {
	return &cli{
		ctx:  ctx,
		log:  log,
		slog: slog,
	}
}
