package s3

import (
	"context"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(options *core.Options) (context.Context, *minio.Client) {
	// TODO get these values from config.yml or from env vars
	ctx := context.Background()
	endpoint := options.S3Endpoint
	accessKeyID := options.S3AccessKey
	secretAccessKey := options.S3SecretAccessKey
	useSSL := options.S3UseSSL

	// Initialize client with
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return ctx, client
}
