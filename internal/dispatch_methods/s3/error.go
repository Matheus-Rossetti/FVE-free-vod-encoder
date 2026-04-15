package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

func (s *S3) rollback(uploadedKeys []minio.ObjectInfo) {
	if len(uploadedKeys) == 0 {
		return
	}

	s.slog.Warn("Starting rollback...")

	// THIS SLEEP IS VERY IMPORTANT!!!
	// If we run the code below, asking for S3 to delete the uploaded files
	// instantly after we cancel the context, sometimes the server gets the
	// delete request midway through saving some files.
	// S3 checks if the file exists to delete it, but since it hasn't been saved yet
	// the server thinks it doesn't exists and doesn't delete it.
	// A couple milliseconds later, the file is saved, behold, an orphan file has been created.
	// So we give some time for the server to actually save all files.
	time.Sleep(time.Second * 4)

	objInfoChan := make(chan minio.ObjectInfo)
	go func() {
		defer close(objInfoChan)
		for _, key := range uploadedKeys {
			objInfoChan <- key
		}
	}()

	rmObjsOptions := minio.RemoveObjectsOptions{}
	errChan := s.client.RemoveObjects(
		context.Background(),
		s.bucket,
		objInfoChan,
		rmObjsOptions,
	)

	for rmErr := range errChan {
		if rmErr.Err != nil {
			s.slog.Error("error removing a file from S3", "file", rmErr.ObjectName, "err", rmErr.Err)
		}
	}

	s.slog.Warn("Rollback finished!")

}
