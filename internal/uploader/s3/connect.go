package s3

import (
	"context"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(options *core.Options) (context.Context, *minio.Client) {

	ctx := context.Background()
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

	return ctx, client
}
