package core

import (
	"log"
	"log/slog"
)

type frevod struct {
	Log     *log.Logger
	Slog    *slog.Logger
	Options *Options
}

func StartFrevod(log *log.Logger, slog *slog.Logger) *frevod {
	return &frevod{
		Log:  log,
		Slog: slog,
	}
}
