package rest

import (
	"context"
	"log"
	"log/slog"
)

type rest struct {
	ctx  context.Context
	log  *log.Logger
	slog *slog.Logger
	port string
}

func NewRest(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	port string,
) *rest {
	return &rest{
		ctx:  ctx,
		log:  log,
		slog: slog,
		port: port,
	}
}
