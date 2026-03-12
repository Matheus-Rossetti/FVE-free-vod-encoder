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
	downloadQueue chan<- *core.Job
}

func NewRest(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	port string,
	downloadQueue chan<- *core.Job,
) *rest {
	return &rest{
		ctx:           ctx,
		port:          port,
		downloadQueue: downloadQueue,
	}
}
