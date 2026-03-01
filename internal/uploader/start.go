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

func Start(id int, uploadQueue <-chan *core.Job, options *core.Options) {

	// START CONNECTIONS
	var ctx context.Context
	var client *minio.Client
	if options.UseS3 {
		ctx, client = s3.Connect(options)
	}

	// CREATE POOL
	uploadPool := make(chan struct{}, options.ConcurrentUploads)
	for range options.ConcurrentUploads {
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

				<-uploadPool // takes a file to upload
				wg.Add(1)

				if options.UseS3 {

					go func() {
						defer wg.Done()
						s3_key := s3.GetKey(job, path)
						absolutePath, _ := filepath.Abs(path)
						// the upload will aways return a struct back to the pool
						s3.UploadFiles(client, s3_key, ctx, absolutePath, uploadPool)
					}()

				}

				return nil
			}) // walkdir finishline

		wg.Wait()

		fmt.Printf("\nFinished uploading!")
		fmt.Printf("\nDeleting ROT files...")
		go DeleteROT(job.UploadJob.Dir)
		fmt.Printf("\nDir %v deleted!", job.UploadJob.Dir)

		// loop
	}
}
