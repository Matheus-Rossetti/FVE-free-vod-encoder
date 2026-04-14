package local

import (
	"os"
)

func (l *local) HandleError(dir string) error {

	err := os.RemoveAll(dir)
	if err != nil {
		l.slog.Error("error deleting a dir after dispatch", "err", err)
	}

	return nil
}
