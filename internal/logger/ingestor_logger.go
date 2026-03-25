package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Ingestor(id int) (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: fmt.Sprintf("<ingestor %v>", id),
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
