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
	dispatchQueue    core.DispatchQueue
	storageProviders []StorageProvider
}

func NewDispatcher(
	ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	options *core.Options,
	dispatchQueue core.DispatchQueue,
	storageProviders []StorageProvider,
) *dispatcher {
	return &dispatcher{
		ctx:              ctx,
		log:              log,
		slog:             slog,
		options:          options,
		dispatchQueue:    dispatchQueue,
		storageProviders: storageProviders,
	}
}
