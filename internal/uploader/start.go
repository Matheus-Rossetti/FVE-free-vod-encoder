package uploader

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

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
			fmt.Println("\nUploading to S3...")

			filepath.WalkDir(
				job.UploadJob.Dir,
				func(path string, entry fs.DirEntry, err error) error {

					if err != nil {
						return err
					}

					if entry.IsDir() {
						return nil
					}

					fmt.Printf("\nTook one from the pool: %v", len(uploadPool))
					<-uploadPool
					fmt.Printf("\nTook one from the pool: %v", len(uploadPool))
					wg.Add(1)
					go func() error {
						defer func() { uploadPool <- struct{}{}; wg.Done() }()
						err = provider.Upload(ctx, job, path)
						if err != nil {
							provider.HandleError()
							return err
						}

						return nil
					}()

					return nil
				}) // walkdir finishline
		}

		wg.Wait()

		fmt.Printf("\nFinished uploading!")
		fmt.Printf("\nDeleting local ROT files...")
		go DeleteROT(job.UploadJob.Dir)
		fmt.Printf("\nDir %v deleted!", job.UploadJob.Dir)
	}

	// After queue closes
	fmt.Printf("\nUploader shuting down...")
}
