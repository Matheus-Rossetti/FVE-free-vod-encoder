package ingestor

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

func (i *ingestor) ingestFromUrl(file *os.File, url string) error {
	i.slog.Info("Downloading...", "from", url)

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
	file.Read(headerBytes)

	// Download directly to disc
	_, err = io.Copy(file, body)
	if err != nil {
		i.filePool <- file
		return fmt.Errorf("%w: %v", ErrReadingBodyIntoFile, err)
	}

	return nil
}

func (i *ingestor) handleUrlError(err error) {

}
