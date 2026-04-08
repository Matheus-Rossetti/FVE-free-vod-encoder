package s3

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

var (
	ErrUploadingFileToS3 = errors.New("failed when uploading a file to S3")
)

func (s *S3) Upload(ctx context.Context, key, filePath string) error {

	// keys built on windows comes with back slashes '\' but S3 uses forward slashes '/'
	key = strings.ReplaceAll(key, `\`, "/")

	_, err := s.Client.FPutObject(
		ctx,
		s.Bucket,
		key,
		filePath,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		s.slog.Error(ErrUploadingFileToS3.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrUploadingFileToS3, err)
	}

	return nil
}
