package s3

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/minio/minio-go/v7"
)

func (s *S3) Upload(ctx context.Context, job *core.Job, path string) error {

	fmt.Printf("\nUploading %v..", path)

	key := GetKey(job, path)
	absolutePath, _ := filepath.Abs(path)

	info, err := s.Client.FPutObject(
		ctx,
		s.Bucket,
		key,
		absolutePath,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		fmt.Printf("\nError uploading: %v", err)
		return fmt.Errorf("Error uploading files: %v", err)
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", key, info.Size)
	return nil
}
