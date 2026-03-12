package main

import (
	"log"
	"log/slog"
)

type frevod struct {
	log  *log.Logger
	slog *slog.Logger
}

func startFrevod(log *log.Logger, slog *slog.Logger) *frevod {
	return &frevod{
		log:  log,
		slog: slog,
	}
}
