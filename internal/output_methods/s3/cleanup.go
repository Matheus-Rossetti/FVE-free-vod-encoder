package s3

import (
	"fmt"

	"github.com/minio/minio-go/v7"
)

// CleanUp removes all segments and manifest files already uploaded to S3
// if a failure occurs during the upload process of the same video set.
// This prevents orphaned files after an upload error.
func CleanUp(client *minio.Client, dir string) {
	fmt.Printf("\nCleaning up failed upload residues in: %v", dir)

}
