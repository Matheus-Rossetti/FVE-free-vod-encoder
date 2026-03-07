package s3

import (
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateClient(options *core.Options) *minio.Client {

	endpoint := options.Upload.S3.S3Endpoint
	accessKeyID := options.Upload.S3.S3AccessKey
	secretAccessKey := options.Upload.S3.S3SecretAccessKey
	useSSL := options.Upload.S3.S3UseSSL

	// Initialize client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
