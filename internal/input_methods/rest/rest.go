package rest

import (
	"context"
	"log"
	"log/slog"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type rest struct {
	ctx           context.Context
	log           *log.Logger
	slog          *slog.Logger
	port          string
	downloadQueue core.IngestQueue
}

func NewRest(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	port string,
	downloadQueue core.IngestQueue,
) *rest {
	return &rest{
		ctx:           ctx,
		log:           log,
		slog:          slog,
		port:          port,
		downloadQueue: downloadQueue,
	}
}
