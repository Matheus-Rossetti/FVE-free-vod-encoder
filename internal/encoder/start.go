package encoder

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrPreparingOutputStorage = errors.New("failed preparing output storage")
	ErrRunningFFmpegCommand   = errors.New("failed when running FFmpeg command")
)

func (e *encoder) Start() {
	for job := range e.encodeQueue {
		func() { // wrap the loop contents in a function so we can defer returning the file to the file pool
			e.slog.Info("Received a job!", "encoding", job.DispatchJob.KeyStarter) // TOOD change for job.Id when job struct gets refactored

			defer func() {
				if job.EncodeJob.DownloadedFile {
					e.filePool <- job.EncodeJob.File
				}
			}()

			videoName := job.EncodeJob.File.Name()
			absoluteVideoPath, err := filepath.Abs(videoName)
			if err != nil {
				e.slog.Error("couldn't get absolute path for video", "video", videoName, "err", err)
				return
			}

			metadata, err := GetMetadataFrom(absoluteVideoPath)
			if err != nil {
				e.slog.Error("couldn't get metadata from a video", "video", videoName, "err", err)
			}

			video := NewVideo(metadata)
			cmd := e.BuildFFmpegCommand(video, absoluteVideoPath)

			outputDir, err := e.createOutputDir()
			if err != nil {
				e.slog.Error(ErrPreparingOutputStorage.Error(), "err", err)
				return
			}

			cmd.Dir = outputDir
			err = cmd.Run()
			if err != nil {
				e.slog.Error(ErrRunningFFmpegCommand.Error(), "err", err)
				os.RemoveAll(outputDir)
				return
			}

			job.DispatchJob.FromDir = outputDir

			e.slog.Info("Finished encoding!")
			e.dispatchQueue <- job
		}()
	}

	// After queue closes
	os.RemoveAll("output")
	e.slog.Info("Shutting down...")
}
