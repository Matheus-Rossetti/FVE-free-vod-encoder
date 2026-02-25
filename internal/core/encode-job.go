package core

import (
	"log"
	"os"
	"path/filepath"
)

type EncodeJob struct {
	AbsoluteVideoPath string
	File              *os.File
}

func NewEncodeJob(file *os.File) EncodeJob {
	absoluteVideoPath, err := filepath.Abs(file.Name())
	if err != nil {
		log.Fatal("Error getting the absolute video path to create encode job", err)
	}

	return EncodeJob{
		AbsoluteVideoPath: absoluteVideoPath,
		File:              file,
	}
}
