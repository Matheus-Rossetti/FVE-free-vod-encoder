package uploader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
)

var (
	ErrGettingAbsolutePath = errors.New("failed getting the absolute path to upload from")
	ErrUploadingFile       = errors.New("Failed when uploading a file")
)

// This function is getting messier by de day, refactor or be ashamed.
func (u *uploader) Start() {

	for job := range u.uploadQueue {
		u.slog.Info(fmt.Sprintf("Received a job! Walking ", job.UploadJob.FromDir), "id", u.id)

		// CREATE UPLOAD POOL
		uploadPool := make(chan struct{}, u.options.Upload.ConcurrentUploads)
		for range u.options.Upload.ConcurrentUploads {
			uploadPool <- struct{}{}
		}

		var wg sync.WaitGroup

		for providerName, provider := range u.storageProviders {
			providerContext, cancel := context.WithCancel(u.ctx)

			var uploadedFiles []string
			var mu sync.Mutex

			filepath.WalkDir(
				job.UploadJob.FromDir,
				func(path string, entry fs.DirEntry, err error) error {

					if providerContext.Err() != nil {
						return providerContext.Err()
					} // providerContext is closed when ctrl + c or an error occurs

					if err != nil || entry.IsDir() {
						return err // if err here is nil, it wont stop .WalkDir
					}

					<-uploadPool

					wg.Add(1)
					go func() error {
						defer func() { uploadPool <- struct{}{}; wg.Done() }()

						key := getKey(job.UploadJob.KeyStarter, path)
						absolutePath, err := filepath.Abs(path)
						if err != nil {
							cancel()
							u.slog.Error(ErrGettingAbsolutePath.Error(), "err", err, "id", u.id)
							return fmt.Errorf("%w: %v", ErrGettingAbsolutePath, err)
						}

						err = provider.Upload(providerContext, key, absolutePath)
						if err != nil {
							cancel() // stops the walkdir
							u.slog.Error(ErrUploadingFile.Error(), "err", err, "id", u.id)
							return fmt.Errorf("%w: %v", ErrUploadingFile, err)
						}

						mu.Lock()
						uploadedFiles = append(uploadedFiles, key)
						mu.Unlock()

						return nil
					}()

					return nil
				})

			wg.Wait() // finishes uploading to one provider before starting another

			if providerContext.Err() != nil { // canceling the original ctx (using ctrl + c) will also cancel the providerContext
				provider.HandleError(uploadedFiles)
			}

			u.slog.Info(fmt.Sprintf("Done upliading to %v", providerName), "dir", job.UploadJob.FromDir, "id", u.id)
		} // provider loop

		go DeleteROT(job.UploadJob.FromDir)
		u.slog.Info("Finished!", "id", u.id)
	}

	// After queue closes
	u.slog.Info("Shutting down...", "id", u.id)
}
