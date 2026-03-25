package logger

import (
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Options() (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "<options>",
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
