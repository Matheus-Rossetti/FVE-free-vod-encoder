package cli

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (c *cli) ListenForInput() {
	c.slog.Info("Frevod is waiting for URIs in your terminal! Use: videoUri key")

	scanner := bufio.NewScanner(os.Stdin)

	var (
		input           string
		parts           []string
		videoUri        string
		videoKeyStarter string
	)

	for {

		if !scanner.Scan() {
			break
		}

		input = scanner.Text()
		parts = strings.Fields(input)

		if len(parts) != 2 {
			c.slog.Error("Need URI and KEY!")
			c.slog.Error("Example: /where/it/is where/to/put")
			continue
		}

		videoUri = parts[0]
		videoKeyStarter = parts[1]

		uriType, err := core.CategorizeUri(videoUri)
		if errors.Is(err, core.ErrUnsupportedUri) {
			c.slog.Error("Unsupported URI!")
			c.slog.Info("Supported URIs starts with: http, https and file.")
			continue
		}

		job := core.NewJob()
		job.DownloadJob.Source = "cli"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = videoUri

		job.UploadJob.KeyStarter = videoKeyStarter

		c.downloadQueue <- job

		time.Sleep(time.Second) // avoid spam
	}
}
