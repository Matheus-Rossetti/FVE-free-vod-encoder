package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (s *S3) Upload(ctx context.Context, client *minio.Client) {

	bucketName := "videos"

	info, err := client.FPutObject(
		ctx,
		bucketName,
		key,
		path,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		return fmt.Errorf("Error uploading files:", err)
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", key, info.Size)
	return nil
}
