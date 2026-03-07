package uploader

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, options *core.Options, uploadQueue <-chan *core.Job, storageProviders ...StorageProvider) {

	for job := range uploadQueue {

		// CREATE UPLOAD POOL
		uploadPool := make(chan struct{}, options.Upload.ConcurrentUploads)
		for range options.Upload.ConcurrentUploads {
			uploadPool <- struct{}{}
		}

		var wg sync.WaitGroup
		filepath.WalkDir(
			job.UploadJob.Dir,
			func(path string, entry fs.DirEntry, err error) error {

				if err != nil {
					return err
				}

				if entry.IsDir() {
					return nil
				}

				for _, provider := range storageProviders {
					provider.Upload(ctx, options, job)
				}

				return nil
			}) // walkdir finishline

		wg.Wait()

		fmt.Printf("\nFinished uploading!")
		fmt.Printf("\nDeleting ROT files...")
		go DeleteROT(job.UploadJob.Dir)
		fmt.Printf("\nDir %v deleted!", job.UploadJob.Dir)
	}

	// After queue closes
	fmt.Printf("\nUploader shuting down...")
}
