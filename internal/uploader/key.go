package uploader

import (
	"path/filepath"
)

func getKey(keyStarter, path string) string {

	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	var s3_key string
	if dir != "." { // dir is "." if path is only the filename
		s3_key = filepath.Join(keyStarter, dir, filename)
	} else {
		s3_key = filepath.Join(keyStarter, filename)
	}

	return s3_key
}
