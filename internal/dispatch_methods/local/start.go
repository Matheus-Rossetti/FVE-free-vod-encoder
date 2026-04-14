package local

import (
	"log"
	"log/slog"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type local struct {
	log        *log.Logger
	slog       *slog.Logger
	storageDir string
}

func Start(log *log.Logger, slog *slog.Logger, options *core.Options) *local {
	slog.Info("Storing locally!", "dir", options.Dispatch.Local.StoreAt)

	return &local{
		log:        log,
		slog:       slog,
		storageDir: options.Dispatch.Local.StoreAt,
	}
}

func (l *local) Name() string {
	return "local"
}
