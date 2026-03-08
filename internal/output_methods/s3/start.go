package s3

import (
	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(options *core.Options) *S3 {
	bucket := options.Upload.S3.S3BucketName
	client := CreateClient(options)

	return &S3{
		Bucket: bucket,
		Client: client,
	}
}
