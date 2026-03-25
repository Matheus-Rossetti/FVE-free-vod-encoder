package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	charmLog "github.com/charmbracelet/log"
)

func Encoder(id int) (*log.Logger, *slog.Logger) {

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: fmt.Sprintf("<encoder %v>", id),
		})

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
