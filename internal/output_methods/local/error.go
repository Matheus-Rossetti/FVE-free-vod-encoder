package local

import (
	"errors"
	"fmt"
	"os"
)

var ErrDeletingFile = errors.New("failed when deleting file")

func (l *local) HandleError(files []string) error {

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			l.slog.Error(ErrDeletingFile.Error(), "file", file)
			return fmt.Errorf("%w: %v", ErrDeletingFile, err)
		}
	}

	return nil
}
