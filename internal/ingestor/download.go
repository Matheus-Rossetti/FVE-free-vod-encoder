package ingestor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func prepareFileForDownload(file *os.File) error {
	if file == nil {
		return fmt.Errorf("%v", "file's closed")
	}

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

	fmt.Printf("\nMIME: %v", MIMEType)
	fmt.Printf("\nfile type: %v", fileType)

	if fileType[0] != "video" {
		return nil, fmt.Errorf("file not a video")
	}

	return buffer, nil
}
