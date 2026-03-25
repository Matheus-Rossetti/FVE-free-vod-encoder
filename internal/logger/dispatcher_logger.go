package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Dispatcher(id int) (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: fmt.Sprintf("<dispatcher %v>", id),
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
