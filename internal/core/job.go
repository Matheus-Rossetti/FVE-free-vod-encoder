package core

import (
	"log"
	"os"
	"path/filepath"
)

type Job struct {
	File              *os.File
	OutputDir         string
	VideoUri          string
	Source            string
	AbsoluteVideoPath string
}

func NewJob(file *os.File, outputDir, videoUri, source string) *Job {
	return &Job{
		File:      file,
		OutputDir: outputDir,
		VideoUri:  videoUri,
		Source:    source,
	}
}

func (j *Job) SetAbsolutePath(path string) {
	AbsoluteVideoPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("Error getting the absolute video path after downloading", err)
	}

	j.AbsoluteVideoPath = AbsoluteVideoPath
}
