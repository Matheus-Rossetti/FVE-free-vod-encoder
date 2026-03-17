package encoder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrPreparingOutputStorage = errors.New("failed preparing output storage")
	ErrGettingAboslutePath    = errors.New("failed getting the absolute path from file name")
	ErrCreatingVideoStruct    = errors.New("failed when collecting video data")
	ErrRunningFFmpegCommand   = errors.New("Failed when running FFmpeg command")
)

func (e *encoder) Start() {
JobLoop:
	for job := range e.encodeQueue {
		defer func() {
			if job.EncodeJob.DownloadedFile {
				e.filePool <- job.EncodeJob.File
			}
		}() // return the file to the pool

		e.slog.Info(fmt.Sprintf("Receive a job! Encoding contents from %v", job.EncodeJob.File.Name()),
			"id", e.id)

		absoluteVideoPath, err := filepath.Abs(job.EncodeJob.File.Name())
		if err != nil {
			e.slog.Error(ErrGettingAboslutePath.Error(), "err", err, "id", e.id)
			continue JobLoop
		}

		video, err := e.NewVideo(absoluteVideoPath)
		if err != nil {
			e.slog.Error(ErrCreatingVideoStruct.Error(), "err", err)
			continue JobLoop
		}

		cmd := e.BuildFFmpegCommand(video)

		outputDir, err := e.createOutputDir()
		if err != nil {
			e.slog.Error(ErrPreparingOutputStorage.Error(), "err", err, "id", e.id)
			continue JobLoop
		}

		cmd.Dir = outputDir
		err = cmd.Run()
		if err != nil {
			e.slog.Error(ErrRunningFFmpegCommand.Error(), "err", err, "id", e.id)
			e.slog.Info(fmt.Sprintf("Deleting %v", outputDir))
			os.RemoveAll(outputDir)
			continue JobLoop
		}

		job.UploadJob.FromDir = outputDir

		e.slog.Info("Finished!", "id", e.id)
		e.uploadQueue <- job
	}

	// After queue closes
	os.RemoveAll("output")
	e.slog.Info("Shutting down...", "id", e.id)
}
