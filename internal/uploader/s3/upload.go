package s3

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
)

func UploadFiles(client *minio.Client, key string, ctx context.Context, relativePath string) {

	// Upload the test file
	// Change the value of filePath if the file is in another location
	contentType := "application/octet-stream"

	bucketName := "videos"

	// Upload the test file with FPutObject
	info, err := client.FPutObject(ctx, bucketName, key, relativePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", key, info.Size)
}
