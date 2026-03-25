package dispatcher

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

	concurrentUploads = 10
)

// This function is getting messier by de day, refactor or be ashamed.
func (d *dispatcher) Start() {
	for job := range d.uploadQueue {
		d.slog.Info("Received a job! Walking dir", "dir", job.UploadJob.FromDir, "id", d.id)

		// CREATE UPLOAD POOL
		uploadPool := make(chan struct{}, concurrentUploads)
		for range concurrentUploads {
			uploadPool <- struct{}{}
		}

		var wg sync.WaitGroup

		for providerName, provider := range d.storageProviders {
			providerContext, cancel := context.WithCancel(d.ctx)

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

						key := getKey(job.UploadJob.KeyStarter, job.UploadJob.FromDir, path)
						absolutePath, err := filepath.Abs(path)
						if err != nil {
							cancel()
							d.slog.Error(ErrGettingAbsolutePath.Error(), "err", err, "id", d.id)
							return fmt.Errorf("%w: %v", ErrGettingAbsolutePath, err)
						}

						err = provider.Upload(providerContext, key, absolutePath)
						if err != nil {
							cancel() // stops the walkdir
							d.slog.Error(ErrUploadingFile.Error(), "err", err, "id", d.id)
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
				err := provider.HandleError(uploadedFiles)
				if err != nil {
					d.slog.Error("Couldn't handle error, sorry :(")
				} else {
					d.slog.Info("Finished cleanup!")
				}
			}

			d.slog.Info("Finished upliading!", "to", providerName, "dir", job.UploadJob.FromDir, "id", d.id)
		} // provider loop

		go d.DeleteROT(job.UploadJob.FromDir)
		d.slog.Info("Finished!", "id", d.id)
	}

	// After queue closes
	d.slog.Info("Shutting down...", "id", d.id)
}
