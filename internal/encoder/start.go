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
	ErrRunningFFmpegCommand   = errors.New("failed when running FFmpeg command")
)

func (e *encoder) Start() {
JobLoop:
	for job := range e.encodeQueue {
		e.slog.Info(fmt.Sprintf("Receive a job! Encoding contents from %v", job.EncodeJob.File.Name()),
			"id", e.id)

		absoluteVideoPath, err := filepath.Abs(job.EncodeJob.File.Name())
		if err != nil {
			e.slog.Error(ErrGettingAboslutePath.Error(), "err", err, "id", e.id)
			if job.EncodeJob.DownloadedFile {
				e.filePool <- job.EncodeJob.File
			}
			continue JobLoop
		}

		video, err := e.NewVideo(absoluteVideoPath)
		if err != nil {
			e.slog.Error(ErrCreatingVideoStruct.Error(), "err", err)
			if job.EncodeJob.DownloadedFile {
				e.filePool <- job.EncodeJob.File
			}
			continue JobLoop
		}

		cmd := e.BuildFFmpegCommand(video)

		outputDir, err := e.createOutputDir()
		if err != nil {
			e.slog.Error(ErrPreparingOutputStorage.Error(), "err", err, "id", e.id)
			if job.EncodeJob.DownloadedFile {
				e.filePool <- job.EncodeJob.File
			}
			continue JobLoop
		}

		cmd.Dir = outputDir
		err = cmd.Run()
		if err != nil {
			e.slog.Error(ErrRunningFFmpegCommand.Error(), "err", err, "id", e.id)
			e.slog.Info(fmt.Sprintf("Deleting %v", outputDir))
			os.RemoveAll(outputDir)
			if job.EncodeJob.DownloadedFile {
				e.filePool <- job.EncodeJob.File
			}
			continue JobLoop
		}

		job.UploadJob.FromDir = outputDir

		e.slog.Info("Finished!", "id", e.id)
		e.uploadQueue <- job

		if job.EncodeJob.DownloadedFile {
			e.filePool <- job.EncodeJob.File
		}
	}

	// After queue closes
	os.RemoveAll("output")
	e.slog.Info("Shutting down...", "id", e.id)
}
