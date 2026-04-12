package ingestor

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrGettingAbsolutePath = errors.New("failed to get absolute path for inputed path")
	ErrFindingLocalFile    = errors.New("failed to find local file")
	ErrLocalFileNotVideo   = errors.New("inputed local file isn't a video")
)

func ingestFromLocal(path string) (*os.File, error) {

	absolutePath, err := getAbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrGettingAbsolutePath, err)
	}

	file, err := checkExistence(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrFindingLocalFile, err)
	}
	defer file.Close()

	err = checkFileForVideo(file)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrLocalFileNotVideo, err)
	}

	return file, nil
}

func getAbsolutePath(path string) (string, error) {
	path = strings.TrimPrefix(path, "file://")

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}

func checkExistence(absolutePath string) (*os.File, error) {
	localFile, err := os.Open(absolutePath)

	if err != nil {
		return nil, err
	}

	return localFile, nil
}

func checkFileForVideo(file *os.File) error {
	buffer, err := io.ReadAll(io.LimitReader(file, 512))
	if err != nil {
		return err
	}

	MIMEType := http.DetectContentType(buffer)
	fileType := strings.Split(MIMEType, "/")

	if fileType[0] != "video" {
		return err
	}

	return nil
}
