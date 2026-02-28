package uploader

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/uploader/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Start(id int, uploadQueue <-chan *core.Job) {

	// TODO get these values from config.yml or from env vars
	ctx := context.Background()
	endpoint := ""
	accessKeyID := ""
	secretAccessKey := ""
	useSSL := true

	// Initialize client with
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
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

				s3_key := s3.GetKey(job, path)
				absolutePath, _ := filepath.Abs(path)
				fmt.Printf("\nAbsolute path for file: %v\n", absolutePath)
				s3.UploadFiles(client, s3_key, ctx, absolutePath)

				fmt.Printf("S3 Key: %v\n", s3_key)

				return nil
			})

		// upload files
	}
}
