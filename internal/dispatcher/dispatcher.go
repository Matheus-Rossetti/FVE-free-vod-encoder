package dispatcher

import (
	"context"
	"log"
	"log/slog"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type dispatcher struct {
	ctx              context.Context
	log              *log.Logger
	slog             *slog.Logger
	options          *core.Options
	id               int
	uploadQueue      <-chan *core.Job
	storageProviders map[string]StorageProvider
}

func NewDispatcher(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	options *core.Options,
	id int,
	uploadQueue <-chan *core.Job,
	storageProviders map[string]StorageProvider,
) *dispatcher {
	return &dispatcher{
		ctx:              ctx,
		log:              log,
		slog:             slog,
		options:          options,
		id:               id,
		uploadQueue:      uploadQueue,
		storageProviders: storageProviders,
	}
}
