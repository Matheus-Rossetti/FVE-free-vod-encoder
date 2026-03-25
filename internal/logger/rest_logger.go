package logger

import (
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Rest() (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "Input/> rest",
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
