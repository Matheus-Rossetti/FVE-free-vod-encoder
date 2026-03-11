package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (s *S3) Upload(ctx context.Context, key, filePath string) error {

	info, err := s.Client.FPutObject(
		ctx,
		s.Bucket,
		key,
		filePath,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		fmt.Printf("\nError uploading: %v", err)
		return fmt.Errorf("Error uploading files: %v", err)
	}

	fmt.Printf("\nSuccessfully uploaded %s of size %d", key, info.Size)
	return nil
}
