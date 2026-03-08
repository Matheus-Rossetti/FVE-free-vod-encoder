package s3

import "github.com/minio/minio-go/v7"

type S3 struct {
	Bucket string
	Client *minio.Client
}
