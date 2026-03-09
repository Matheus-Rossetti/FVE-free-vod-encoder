package s3

import (
	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
)

type S3 struct {
	Bucket string
	Client *minio.Client
}

func Start(options *core.Options) *S3 {
	bucket := options.Upload.S3.S3BucketName
	client := CreateClient(options)

	return &S3{
		Bucket: bucket,
		Client: client,
	}
}
