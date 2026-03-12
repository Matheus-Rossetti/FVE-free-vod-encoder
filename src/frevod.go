package main

import (
	"log"
	"log/slog"
)

type frevod struct {
	log  *log.Logger
	slog *slog.Logger
}

func StartFrevod(log *log.Logger, slog *slog.Logger) *frevod {
	return &frevod{
		log:  log,
		slog: slog,
	}
}
