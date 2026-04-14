package local

import (
	"context"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (l *local) HandleError(ctx context.Context, job core.DispatchJob) error {

	err := os.RemoveAll(job.FromDir)
	if err != nil {
		l.slog.Error("error deleting a dir after dispatch", "err", err)
	}

	return nil
}
