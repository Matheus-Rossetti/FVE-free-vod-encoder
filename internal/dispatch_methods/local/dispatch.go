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

	baseDirName := path.Join(l.storageDir, job.KeyStarter)

	counter := 0
	newDirName := baseDirName

	for { // if dir already exists, we add {counter} to the name
		_, err := os.Stat(newDirName)
		if errors.Is(err, fs.ErrNotExist) {
			break
		}

		oldName := newDirName
		counter++

		newDirName = fmt.Sprintf("%s%d", baseDirName, counter)
		l.slog.Warn("This dir already exists!", "dir", oldName, "trying as", newDirName)
	}

	err := os.Rename(job.FromDir, newDirName)
	if err != nil {
		l.HandleError(job.FromDir) // TODO fall back to creating a new dir and copying the files
		return fmt.Errorf("failed to rename dir from %v to %v | %v", job.FromDir, newDirName, err)
	}

	l.slog.Info("Saved!", "here", newDirName)
	return nil
}
