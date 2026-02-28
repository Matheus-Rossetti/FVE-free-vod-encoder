package s3

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
)

func UploadFiles(client *minio.Client, key string, ctx context.Context, relativePath string, uploadSlot chan<- struct{}) {

	contentType := "video/mp4"

	bucketName := "videos"

	// Upload the test file with FPutObject
	info, err := client.FPutObject(ctx, bucketName, key, relativePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", key, info.Size)
	uploadSlot <- struct{}{}
}
