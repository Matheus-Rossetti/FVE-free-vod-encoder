package s3

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

var ErrRemovingObj = errors.New("Failed to remove and object from S3")

func (s *S3) HandleError(uploadedFiles []string) error {

	s.slog.Warn("Oops.. Starting cleanup...")

	if len(uploadedFiles) < 1 {
		s.slog.Warn("Nothing to clean, no files were uploaded!")
		return nil
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
			s.slog.Error(ErrRemovingObj.Error(), "object", rmErr.ObjectName, "err", rmErr.Err)
			return fmt.Errorf("%w: %v", ErrRemovingObj, rmErr.Err)
		}
	}

	return nil
}
