package uploader

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/uploader/s3"
)

func Start(id int, uploadQueue <-chan *core.Job, options *core.Options) {

	ctx, client := s3.Connect(options)
	simultaneousUploads := 10
	uploadSlot := make(chan struct{}, simultaneousUploads)
	for range simultaneousUploads {
		uploadSlot <- struct{}{}
	}

	for job := range uploadQueue {

		filepath.WalkDir(
			job.UploadJob.DirToUploadFrom,
			func(path string, entry fs.DirEntry, err error) error {

				if err != nil {
					return err
				}

				if entry.IsDir() {
					return nil
				}

				<-uploadSlot // takes one

				s3_key := s3.GetKey(job, path)
				absolutePath, _ := filepath.Abs(path)
				go s3.UploadFiles(client, s3_key, ctx, absolutePath, uploadSlot)

				return nil
			})

		fmt.Printf("Finished uploading!")

		// delete files after upload
	}
}
