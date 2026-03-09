package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func (s *S3) HandleError(uploadedFiles []string) {

	fmt.Printf("\nStarting cleanup")

	if len(uploadedFiles) < 1 {
		fmt.Println("No files were uploaded.")
		return
	}

	objInfoChan := make(chan minio.ObjectInfo)

	go func() {
		defer close(objInfoChan)

		for _, key := range uploadedFiles {
			objInfoChan <- minio.ObjectInfo{
				Key: key,
			}
		}
	}()

	rmObjsOptions := minio.RemoveObjectsOptions{}

	errChan := s.Client.RemoveObjects(context.Background(), s.Bucket, objInfoChan, rmObjsOptions)

	for rmErr := range errChan {
		if rmErr.Err != nil {
			fmt.Printf("\nFailed to remove object '%s': %v", rmErr.ObjectName, rmErr.Err)
		}
	}

	fmt.Printf("\n\n\nCleanup process finished!")
}
