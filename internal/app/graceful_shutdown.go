package app

import (
	"os"
	"path/filepath"
)

func (a *app) CloseAndDeleteStorageFiles(filePool <-chan *os.File) {
	for file := range filePool {
		file.Close()
		absoluteFilePath, _ := filepath.Abs(file.Name())
		os.Remove(absoluteFilePath)
		a.slog.Warn("Closed and Deleted", "file", file.Name())
	}
}
