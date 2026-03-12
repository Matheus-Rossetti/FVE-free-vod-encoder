package options

import (
	"log"
	"log/slog"
)

type options struct {
	log  *log.Logger
	slog *slog.Logger
}

func NewOptions(log *log.Logger, slog *slog.Logger) *options {
	return &options{
		log:  log,
		slog: slog,
	}
}
