package s3

import (
	"log"
	"log/slog"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	log    *log.Logger
	slog   *slog.Logger
	Bucket string
	Client *minio.Client
}

func Start(log *log.Logger, slog *slog.Logger, options *core.Options) *S3 {

	bucket := options.Upload.S3.BucketName
	endpoint := options.Upload.S3.Endpoint
	accessKeyID := options.Upload.S3.AccessKey
	secretAccessKey := options.Upload.S3.SecretAccessKey
	useSSL := options.Upload.S3.UseSSL

	// Minio complains if the endpoint includes http://
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	// Initialize client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to start S3 upload method. err=%v", err)
	}

	slog.Info("Storing on S3!", "bucket", bucket)

	return &S3{
		log:    log,
		slog:   slog,
		Bucket: bucket,
		Client: client,
	}
}
