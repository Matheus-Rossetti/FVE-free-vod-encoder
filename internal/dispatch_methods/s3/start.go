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
	bucket string
	client *minio.Client
}

func Start(log *log.Logger, slog *slog.Logger, options *core.Options) *S3 {

	bucket := options.Dispatch.S3.Bucket
	endpoint := options.Dispatch.S3.Endpoint
	accessKeyID := options.Dispatch.S3.AccessKey
	secretAccessKey := options.Dispatch.S3.SecretAccessKey
	useSSL := options.Dispatch.S3.UseSSL

	// Minio package complains if the endpoint includes http:// or https://
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	// Initialize client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to start S3 Dispatch method. err=%v", err)
	}

	slog.Info("Storing on S3!", "bucket", bucket)

	return &S3{
		log:    log,
		slog:   slog,
		bucket: bucket,
		client: client,
	}
}

func (s *S3) Name() string {
	return "S3"
}
