package encoder

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrCreatingOutputDir = errors.New("failed to create the output/ dir to store encoded segments")
	ErrCreatingTempDir   = errors.New("failed to create a temporary dir")
)

func (e *encoder) createOutputDir() (string, error) {
	err := os.MkdirAll("output", 0700)
	if err != nil {
		e.slog.Error(ErrCreatingOutputDir.Error(), "err", err)
		return "", fmt.Errorf("%w: %v", ErrCreatingOutputDir, err)
	}

	// The return value here is "output/video-{random-string}"
	dirName, err := os.MkdirTemp("output", "video-*")
	if err != nil {
		e.slog.Error(ErrCreatingTempDir.Error(), "err", err)
		return "", fmt.Errorf("%w: %v", ErrCreatingTempDir, err)
	}

	return dirName, nil
}
