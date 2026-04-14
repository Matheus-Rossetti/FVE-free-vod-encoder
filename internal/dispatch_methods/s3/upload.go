package s3

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
)

var (
	ErrUploadingFileToS3 = errors.New("failed when uploading a file to S3")
)

func (s *S3) Upload(ctx context.Context, job core.DispatchJob) error {

	uploadPool := make(chan struct{}, 10)

	var uploadedFiles []string
	var mu sync.Mutex

	filepath.WalkDir(
		job.UploadJob.FromDir,
		func(currentEntryPath string, entry fs.DirEntry, err error) error {

			if err != nil || entry.IsDir() {
				return err // if isDir == true, err will be nil and walkdir will continue
			}

			uploadPool <- struct{}{}

			wg.Add(1)
			go func() error {
				defer func() { <-uploadPool; wg.Done() }()

				key := getKey(job.UploadJob.KeyStarter, job.UploadJob.FromDir, currentEntryPath)
				absolutePath, err := filepath.Abs(path)
				if err != nil {
					d.slog.Error(ErrGettingAbsolutePath.Error(), "err", err)
					return fmt.Errorf("%w: %v", ErrGettingAbsolutePath, err)
				}

				err = provider.Dispatch(d.ctx, key, absolutePath)
				if err != nil {
					d.slog.Error(ErrUploadingFile.Error(), "err", err)
					return fmt.Errorf("%w: %v", ErrUploadingFile, err)
				}

				mu.Lock()
				uploadedFiles = append(uploadedFiles, key)
				mu.Unlock()

				return nil
			}()

			return nil
		})

	// keys built on windows comes with back slashes '\' but S3 uses forward slashes '/'
	key = strings.ReplaceAll(key, `\`, `/`)

	_, err := s.Client.FPutObject(
		ctx,
		s.Bucket,
		key,
		filePath,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		s.slog.Error(ErrUploadingFileToS3.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrUploadingFileToS3, err)
	}

	return nil
}
