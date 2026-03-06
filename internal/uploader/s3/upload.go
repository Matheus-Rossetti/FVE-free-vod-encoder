package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func UploadFiles(ctx context.Context, client *minio.Client, key string, path string) error {

	bucketName := "videos"

	info, err := client.FPutObject(
		ctx,
		bucketName,
		key,
		path,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		fmt.Println("Error uploading files:", err)
		return err
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", key, info.Size)
	return nil
}
