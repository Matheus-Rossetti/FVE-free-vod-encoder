package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

func (s *S3) HandleError(uploadedFiles []string) {

	fmt.Printf("\nStarting cleanup")

	if len(uploadedFiles) < 1 {
		fmt.Println("No files were uploaded.")
		return
	}

	// THIS SLEEP IS VERY IMPORTANT!!!
	// If we run the code below, asking for S3 to delete the uploaded files
	// instantly after we cancel the context, sometimes the server get the
	// delete request midway through saving some files.
	// S3 checks if the file exists to delete it, but since it hasn't been saved yet
	// the server thinks it doesn't exists and doesn't delete it.
	// A couple milliseconds later, the file is saved, behold, an orphan file has been created.
	// So we give some time for the server to actually save all files.
	time.Sleep(time.Second * 3)

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
