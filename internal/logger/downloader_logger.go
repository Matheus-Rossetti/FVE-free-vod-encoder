package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/charmbracelet/lipgloss"
	charmLog "github.com/charmbracelet/log"
)

func Downloader() (*log.Logger, *slog.Logger) {

	styles := charmLog.DefaultStyles()
	styles.Prefix = lipgloss.NewStyle().Bold(true).Faint(false).Foreground(lipgloss.Color("#5B2A86"))

	charmLogger := charmLog.NewWithOptions(
		os.Stderr,
		charmLog.Options{
			Prefix: "<downloader>",
		})

	charmLogger.SetStyles(styles)

	log := charmLogger.StandardLog()
	slog := slog.New(charmLogger)
	return log, slog
}
