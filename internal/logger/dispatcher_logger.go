package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/charmbracelet/lipgloss"
	charmLog "github.com/charmbracelet/log"
)

func Dispatcher(id int) (*log.Logger, *slog.Logger) {

	styles := charmLog.DefaultStyles()
	styles.Prefix = lipgloss.NewStyle().Bold(true).Faint(false).Foreground(lipgloss.Color("#ff009d"))

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: fmt.Sprintf("<dispatcher %v>", id),
		})

	charmLogger.SetStyles(styles)

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
