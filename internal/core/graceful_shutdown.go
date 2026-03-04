package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func CloseAndDeleteStorageFiles(filePool <-chan *os.File) {
	for file := range filePool {
		file.Close()
		absoluteFilePath, _ := filepath.Abs(file.Name())
		os.Remove(absoluteFilePath)
		fmt.Printf("Closed and Deleted %v\n", file.Name())
	}
}
