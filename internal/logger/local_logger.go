package logger

import (
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Local() (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "<local>",
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
