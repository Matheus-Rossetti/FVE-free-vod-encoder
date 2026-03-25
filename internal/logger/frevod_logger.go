package logger

import (
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Frevod() (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "<frevod>",
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
