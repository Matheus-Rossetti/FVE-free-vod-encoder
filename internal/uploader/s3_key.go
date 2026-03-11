package uploader

import (
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func getKey(job *core.Job, originFilePath string) string {
	relativePath, _ := filepath.Rel(job.UploadJob.Dir, originFilePath)

	dir := filepath.Dir(relativePath)
	filename := filepath.Base(relativePath)

	var s3_key string
	if dir != "." {
		s3_key = filepath.Join(job.UploadJob.S3KeyStarter, dir, filename)
	} else {
		s3_key = filepath.Join(job.UploadJob.S3KeyStarter, filename)
	}

	return s3_key
}
