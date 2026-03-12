package downloader

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

type FileType int

const (
	unknown FileType = iota
	notVideo
	video
)

func (d *downloader) downloadToFile(url string, file *os.File) error {

	// Check if it's a video
	fyleType, err := d.checkFileType(url)
	if err != nil {
		return err
	}
	if fyleType != video {
		d.slog.Error(ErrFileNotVideo.Error())
		return ErrFileNotVideo
	}

	// Create request
	requestCtx, cancel := context.WithTimeout(d.ctx, time.Minute*15)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, "GET", url, nil)
	if err != nil {
		d.slog.Error(ErrFailedRequestCreation.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrFailedRequestCreation, err)
	}

	// Execute request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		d.slog.Error(ErrFailedRequestExecution.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrFailedRequestExecution, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		d.slog.Error(ErrResponseStatusNotOk.Error(), "status", response.Status)
		return fmt.Errorf("%w: status %v", ErrResponseStatusNotOk, response.Status)
	}

	// Download directly to file on disc
	_, err = io.Copy(file, response.Body)
	if err != nil {
		d.slog.Error(ErrReadingBodyIntoFile.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrReadingBodyIntoFile, err)
	}

	return nil
}

// Downloads the first 512 bytes (or less) and check fileType
func (d *downloader) checkFileType(url string) (FileType, error) {

	requestContext, cancel := context.WithTimeout(d.ctx, time.Second*30)
	defer cancel()

	request, err := http.NewRequestWithContext(requestContext, "GET", url, nil)
	if err != nil {
		d.slog.Error(ErrFailedRequestCreation.Error(), "err", err)
		return unknown, fmt.Errorf("%w: %v", ErrFailedRequestCreation, err)
	}

	request.Header.Set("Range", "bytes=0-511")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		d.slog.Error(ErrFailedRequestExecution.Error(), "err", err)
		return unknown, fmt.Errorf("%w: %v", ErrFailedRequestExecution, err)
	}
	defer response.Body.Close()

	// The code below garantees that if the server ignores the Range in the Header and
	// sends the full file, we cut the connection after downloading 512 bytes.
	// This happens becuase of how http.Response.Body works.
	// Hover ".Body" and check the first few lines in the docs, its awesome.
	buffer := make([]byte, 512)
	bytesRead, err := io.ReadFull(response.Body, buffer)
	if err != nil {
		d.slog.Error(ErrReadingBodyIntoBuffer.Error(), "err", err)
		return unknown, fmt.Errorf("%w: %v", ErrReadingBodyIntoBuffer, err)
	}

	// Need to pass :bytesRead in case less than 512 bytes were read into buffer
	MIMEType := http.DetectContentType(buffer[:bytesRead])

	fileType := strings.Split(MIMEType, "/")

	if fileType[0] != "video" {
		return notVideo, nil
	}

	return video, nil
}

func (d *downloader) prepareFileForDownload(file *os.File) error {

	err := file.Truncate(0)
	if err != nil {
		d.slog.Error(ErrTruncatingFile.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrTruncatingFile, err)
	}

	// Point to file head, otherwise might start download somewhere else
	_, err = file.Seek(0, 0)
	if err != nil {
		d.slog.Error(ErrTruncatingFile.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrPointingToFileHead, err)
	}

	return nil
}
