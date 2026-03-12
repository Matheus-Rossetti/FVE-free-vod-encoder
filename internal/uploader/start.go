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
func Start(ctx context.Context, options *core.Options, id int, uploadQueue <-chan *core.Job, storageProviders map[string]StorageProvider) {

	for job := range uploadQueue {
		fmt.Printf("\nUploader %v received %v to upload", id, job.UploadJob.FromDir)

		// CREATE UPLOAD POOL
		uploadPool := make(chan struct{}, options.Upload.ConcurrentUploads)
		for range options.Upload.ConcurrentUploads {
			uploadPool <- struct{}{}
		}

		var wg sync.WaitGroup

		for name, provider := range storageProviders {
			fmt.Printf("Uploading to %v", name)

			providerContext, cancel := context.WithCancel(ctx)

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
						key := getKey(job, path)
						absolutePath, _ := filepath.Abs(path)
						err := provider.Upload(providerContext, key, absolutePath)
						if err != nil {
							fmt.Printf("\nGot an error when Uploading file: %v", err)
							cancel() // stops the walkdir
							return err
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

		fmt.Printf("\nFinished upload job for %v!", job.UploadJob.FromDir)
		fmt.Printf("\nDeleting local ROT files...")
		go DeleteROT(job.UploadJob.FromDir)
		fmt.Printf("\nDir %v deleted!", job.UploadJob.FromDir)
	}

	// After queue closes
	fmt.Printf("\nUploader %v shuting down...", id)
}
