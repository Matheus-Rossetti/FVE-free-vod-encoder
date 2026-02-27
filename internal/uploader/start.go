package uploader

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/uploader/s3"
)

func Start(id int, uploadQueue <-chan *core.Job) {

	// will need base s3 path to upload to, maybe get from input method
	// eg: channel-123/crazy video number 6/
	// then after it I'll put streams_x and the normal and master playlilsts
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

				s3_key := s3.GetKey(job, path)

				fmt.Printf("S3 Key: %v\n", s3_key)

				return nil
			})
	}
}
