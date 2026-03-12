package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/charmbracelet/lipgloss"
	charmLog "github.com/charmbracelet/log"
)

func Cli() (*log.Logger, *slog.Logger) {

	styles := charmLog.DefaultStyles()
	styles.Prefix = lipgloss.NewStyle().Bold(true).Faint(false).Foreground(lipgloss.Color("#FF7F11"))

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "<cli>",
		})

	charmLogger.SetStyles(styles)

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
