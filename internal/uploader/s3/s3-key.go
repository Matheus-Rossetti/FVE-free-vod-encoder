package s3

import (
	"fmt"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func GetKey(job *core.Job, filePath string) string {
	relativePath, _ := filepath.Rel(job.UploadJob.Dir, filePath)

	dir := filepath.Dir(relativePath)
	filename := filepath.Base((relativePath))

	var s3_key string
	if dir != "." {
		s3_key = fmt.Sprintf("%v%v/%v", job.UploadJob.S3KeyStarter, dir, filename)
	} else {
		s3_key = fmt.Sprintf("%v%v", job.UploadJob.S3KeyStarter, filename)
	}

	return s3_key
}
