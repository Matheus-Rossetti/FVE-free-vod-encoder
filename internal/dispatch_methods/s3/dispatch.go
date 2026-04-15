package s3

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
)

func (s *S3) Dispatch(ctx context.Context, job core.DispatchJob) error {

	walkContext, cancelWalk := context.WithCancel(ctx)
	defer cancelWalk()

	uploadSlot := make(chan struct{}, 10)
	var uploadedKeys []minio.ObjectInfo

	mux := sync.Mutex{}
	wg := sync.WaitGroup{}

	filepath.WalkDir(
		job.FromDir,
		func(currentEntryPath string, entry fs.DirEntry, err error) error {
			uploadSlot <- struct{}{}

			if walkContext.Err() != nil {
				return fmt.Errorf("context canceled")
			}

			if err != nil || entry.IsDir() {
				return err // if isDir == true, err will be nil and walkdir will continue
			}

			wg.Go(func() {
				defer func() { <-uploadSlot }()

				key := getKey(job.KeyStarter, job.FromDir, currentEntryPath)
				absolutePath, err := filepath.Abs(currentEntryPath)
				if err != nil {
					s.slog.Error("failed to get absolute path", "path", currentEntryPath, "err", err)
					cancelWalk()
					return
				}

				err = s.uploadFile(walkContext, key, absolutePath)
				if err != nil {
					s.slog.Error("upload failed", "file", currentEntryPath, "err", err)
					cancelWalk()
					// do not return here
				}

				// this could be a goroutine with a channel that updates the slice
				// but this is simpler :)
				mux.Lock()
				uploadedKeys = append(uploadedKeys, minio.ObjectInfo{Key: key})
				mux.Unlock()
			})

			return nil
		})

	wg.Wait() // finish uploading all files before returning
	if walkContext.Err() != nil {
		s.rollback(uploadedKeys)
	}

	return nil
}

func (s *S3) uploadFile(ctx context.Context, key string, filePath string) error {

	// 30sec per segment or manifest
	uploadContext, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, err := s.client.FPutObject(
		uploadContext,
		s.bucket,
		key,
		filePath,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		s.slog.Error("failed to upload file to s3", "file", key, "err", err)
	}

	return err
}
