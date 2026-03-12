package app

import (
	"log"
	"log/slog"
)

type app struct {
	log  *log.Logger
	slog *slog.Logger
}

func NewApp(log *log.Logger, slog *slog.Logger) *app {
	return &app{
		log:  log,
		slog: slog,
	}
}
