package local

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (l *local) Dispatch(ctx context.Context, job core.DispatchJob) error {

	newDirName := path.Join(l.storageDir, job.KeyStarter)
	counter := 1

	for { // if dir already exists, we add {counter} to the name
		_, err := os.Stat(newDirName)
		if errors.Is(err, fs.ErrNotExist) {
			l.slog.Warn("This dir already exists!", "dir", job.FromDir, "Trying as", newDirName)
			break
		} else {
			newDirName = fmt.Sprintf("%s%d", newDirName, counter)
		}
	}

	err := os.Rename(job.FromDir, newDirName)
	if err != nil {
		l.HandleError(job.FromDir) // TODO fall back to creating a new dir and copying the files
		return fmt.Errorf("failed to rename dir from %v to %v | %v", job.FromDir, newDirName, err)
	}

	l.slog.Info("Saved!", "here", newDirName)
	return nil
}
