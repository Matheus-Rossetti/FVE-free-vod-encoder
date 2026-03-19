package uploader

import (
	"path/filepath"
)

func getKey(keyStarter, basePath, targPath string) string {

	relativePath, _ := filepath.Rel(basePath, targPath)

	dir := filepath.Dir(relativePath)
	filename := filepath.Base(relativePath)

	var s3_key string
	if dir != "." { // dir is "." if path is only the filename
		s3_key = filepath.Join(keyStarter, dir, filename)
	} else {
		s3_key = filepath.Join(keyStarter, filename)
	}

	return s3_key
}
