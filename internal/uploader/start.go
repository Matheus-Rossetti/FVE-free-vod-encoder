package uploader

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/uploader/s3"
	"github.com/minio/minio-go/v7"
)

func Start(ctx context.Context, options *core.Options, uploadQueue <-chan *core.Job) {

	// START CONNECTIONS
	var client *minio.Client
	if options.Upload.S3.Use {
		client = s3.Connect(options)
	}

	// CREATE POOL
	uploadPool := make(chan struct{}, options.Upload.ConcurrentUploads)
	for range options.Upload.ConcurrentUploads {
		uploadPool <- struct{}{}
	}

	for job := range uploadQueue {
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

				<-uploadPool

				if options.Upload.S3.Use {
					wg.Add(1)
					go func() {
						s3_key := s3.GetKey(job, path)
						absolutePath, _ := filepath.Abs(path)
						err := s3.UploadFiles(ctx, client, s3_key, absolutePath)
						if err != nil {

						}
						uploadPool <- struct{}{}
						wg.Done()
					}()

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
