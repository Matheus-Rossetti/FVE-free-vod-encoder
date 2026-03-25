package ingestor

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrGettingAbsolutePath = errors.New("failed to get absolute path for inputed path")
	ErrFindingLocalFile    = errors.New("failed to find local file")
)

func (d *ingestor) validateFileFromPath(path string) (*os.File, error) {

	// Format path
	path = strings.TrimPrefix(path, "file://")
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		d.slog.Error(ErrGettingAbsolutePath.Error(), "inputed path", path, "err", err)
		return nil, fmt.Errorf("%w: %v", ErrGettingAbsolutePath, err)
	}

	// Check if file exists
	localFile, err := os.Open(absolutePath)
	if err != nil {
		d.slog.Error(ErrFindingLocalFile.Error(), "err", err)
		return nil, fmt.Errorf("%w: %v", ErrFindingLocalFile, err)
	}
	defer localFile.Close()

	buffer := make([]byte, 512)
	bytesRead, err := localFile.Read(buffer)
	if err != nil {
		d.slog.Error(ErrFindingLocalFile.Error(), "err", err)
		return nil, fmt.Errorf("%w: %v", ErrFindingLocalFile, err)
	}

	// check if it's video
	MIMEtype := http.DetectContentType(buffer[:bytesRead])
	fileType := strings.Split(MIMEtype, "/")
	if fileType[0] != "video" {
		d.slog.Error(ErrFileNotVideo.Error(), "file type", fileType[0])
		return nil, ErrFileNotVideo
	}

	return localFile, nil
}
