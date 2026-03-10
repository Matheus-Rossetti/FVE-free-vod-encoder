package uploader

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

// This function is getting messier by de day, refactor or be ashamed.
func Start(ctx context.Context, options *core.Options, uploadQueue <-chan *core.Job, storageProviders map[string]StorageProvider) {

	for job := range uploadQueue {
		fmt.Printf("\nUploader received %v to upload", job.UploadJob.Dir)

		// CREATE UPLOAD POOL
		uploadPool := make(chan struct{}, options.Upload.ConcurrentUploads)
		for range options.Upload.ConcurrentUploads {
			uploadPool <- struct{}{}
		}

		var wg sync.WaitGroup

		for _, provider := range storageProviders {

			providerContext, cancel := context.WithCancel(ctx)

			var uploadedFiles []string
			var mu sync.Mutex

			filepath.WalkDir(
				job.UploadJob.Dir,
				func(path string, entry fs.DirEntry, err error) error {

					if providerContext.Err() != nil {
						return providerContext.Err()
					}

					if err != nil || entry.IsDir() {
						return err // if err here is nil, it wont stop .WalkDir
					}

					<-uploadPool

					wg.Add(1)
					go func() error {
						defer func() { uploadPool <- struct{}{}; wg.Done() }()
						key, err := provider.Upload(providerContext, job, path)
						if err != nil {
							cancel() // stops the walkdir
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
				fmt.Printf("\n uploadedFiles: %v", uploadedFiles)
				provider.HandleError(uploadedFiles)
			}

		} // provider loop

		fmt.Printf("\nFinished upload job for %v!", job.UploadJob.Dir)
		fmt.Printf("\nDeleting local ROT files...")
		go DeleteROT(job.UploadJob.Dir)
		fmt.Printf("\nDir %v deleted!", job.UploadJob.Dir)
	}

	// After queue closes
	fmt.Printf("\nUploader shuting down...")
}
