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
	ErrPreparingFileforDownload = errors.New("failed to prepare file for download")
	ErrReadingBodyIntoFile      = errors.New("failed to read response body stream into file")
	ErrReadingBodyIntoBuffer    = errors.New("failed to read response body stream into buffer")
	ErrResponseStatusNotOk      = errors.New("request returned a status different from 200")
	ErrFailedRequestCreation    = errors.New("failed to create the request to download video")
	ErrFailedRequestExecution   = errors.New("failed to execute the request to download video")
	ErrTruncatingFile           = errors.New("failed to truncate file before downloading")
	ErrPointingToFileHead       = errors.New("failed to point to file's head")
	ErrFileNotVideo             = errors.New("provided url doesn't download a video file")
)

func (i *ingestor) ingestFromUrl(file *os.File, url string) error {

	err := prepareFileForDownload(file)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPreparingFileforDownload, err)
	}

	// 20 minutes to download a video is enought, right?
	requestCtx, cancel := context.WithTimeout(i.ctx, time.Minute*20)
	defer cancel()

	body, err := makeRequest(url, requestCtx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedRequestCreation, err)
	}
	defer body.Close()

	headerBytes, err := checkBodyForVideo(body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFileNotVideo, err)
	}
	file.Write(headerBytes)

	// Download directly to disc
	_, err = io.Copy(file, body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrReadingBodyIntoFile, err)
	}

	return nil
}

func prepareFileForDownload(file *os.File) error {

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

func makeRequest(url string, ctx context.Context) (io.ReadCloser, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRequestCreation, err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRequestExecution, err)
	}
	// DO NOT CLOSE THE BODY HERE, we'll return it

	// TODO add checks for specific status codes
	status := response.StatusCode
	if status < 200 || status > 300 {
		response.Body.Close()
		return nil, fmt.Errorf("%w: status %v", ErrResponseStatusNotOk, response.Status)
	}

	return response.Body, nil
}

func checkBodyForVideo(body io.ReadCloser) ([]byte, error) {

	buffer, err := io.ReadAll(io.LimitReader(body, 512))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingBodyIntoBuffer, err)
	}

	MIMEType := http.DetectContentType(buffer)
	fileType := strings.Split(MIMEType, "/")

	if fileType[0] != "video" {
		return nil, fmt.Errorf("%w: %v", ErrFileNotVideo, err)
	}

	return buffer, nil
}
