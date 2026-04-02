package ingestor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrFileNotVideo           = errors.New("provided url doesn't download a video file")
	ErrFailedRequestCreation  = errors.New("failed to create the request to download video")
	ErrFailedRequestExecution = errors.New("failed to execute the request to download video")
	ErrResponseStatusNotOk    = errors.New("request returned a status different from 200")
	ErrReadingBodyIntoFile    = errors.New("failed to read response body stream into file")
	ErrReadingBodyIntoBuffer  = errors.New("failed to read response body stream into buffer")
	ErrTruncatingFile         = errors.New("failed to truncate file before downloading")
	ErrPointingToFileHead     = errors.New("failed to point to file's head")
)

func (i *ingestor) prepareFileForDownload(file *os.File) error {

	err := file.Truncate(0)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTruncatingFile, err)
	}

	// Point to the head of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPointingToFileHead, err)
	}

	return nil
}

// Downloads the first 512 bytes (or less) and check fileType
func (i *ingestor) checkForVideo(bytes []byte) error {

	// Need to pass :bytesRead in case less than 512 bytes were read into buffer
	MIMEType := http.DetectContentType(bytes)

	fileType := strings.Split(MIMEType, "/")

	if fileType[0] != "video" {
		return fmt.Errorf("file not a video")
	}

	return nil
}

func (i *ingestor) downloadToFile(url string, file *os.File) error {

	// Create request
	requestCtx, cancel := context.WithTimeout(i.ctx, time.Minute*20)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, "GET", url, nil)
	if err != nil {
		i.slog.Error(ErrFailedRequestCreation.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrFailedRequestCreation, err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		i.slog.Error(ErrFailedRequestExecution.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrFailedRequestExecution, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		i.slog.Error(ErrResponseStatusNotOk.Error(), "status", response.Status)
		return fmt.Errorf("%w: status %v", ErrResponseStatusNotOk, response.Status)
	}

	// The code below garantees that if the server ignores the Range in the Header and
	// sends the full file, we cut the connection after downloading 512 bytes.
	// This happens becuase of how http.Response.Body works.
	// Hover ".Body" and check the first few lines in the docs, its awesome.
	buffer := make([]byte, 512)
	bytesRead, err := io.ReadFull(response.Body, buffer)
	if err != nil {
		i.slog.Error(ErrReadingBodyIntoBuffer.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrReadingBodyIntoBuffer, err)
	}

	err = i.checkForVideo(buffer[:bytesRead])
	if err != nil {
		return err
	}
	file.Write(buffer[:bytesRead])

	// Download directly to file on disc
	_, err = io.Copy(file, response.Body)
	if err != nil {
		i.slog.Error(ErrReadingBodyIntoFile.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrReadingBodyIntoFile, err)
	}

	return nil
}
