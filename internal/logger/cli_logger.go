package logger

import (
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Cli() (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "Input/> cli",
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
