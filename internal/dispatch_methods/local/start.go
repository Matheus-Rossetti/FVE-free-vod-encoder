package local

import (
	"log"
	"log/slog"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type local struct {
	log        *log.Logger
	slog       *slog.Logger
	storageDir string
}

func Start(log *log.Logger, slog *slog.Logger, options *core.Options) *local {

	err := os.MkdirAll(options.Dispatch.Local.StoreAt, 0700)
	if err != nil {
		log.Fatalf("Couldn't create dir to store videos: %v", err)
	}

	slog.Info("Storing locally!", "dir", options.Dispatch.Local.StoreAt)

	return &local{
		log:        log,
		slog:       slog,
		storageDir: options.Dispatch.Local.StoreAt,
	}
}
