package s3

import (
	"path/filepath"
)

func getKey(keyStarter, basePath, targPath string) string {

	// I had a hard time with .Rel, so here's an explanation:
	// basePath = output/video-1234/
	// targPath = output/video-1234/manifest.mp4 or output/video-1234/p480/segment01.ts
	// relativePath = manifest.mp4 or p480/segment01.ts
	relativePath, _ := filepath.Rel(basePath, targPath)

	dir := filepath.Dir(relativePath)
	filename := filepath.Base(relativePath)

	var key string
	if dir != "." { // dir is "." if path is only the filename
		key = filepath.Join(keyStarter, dir, filename)
	} else {
		key = filepath.Join(keyStarter, filename)
	}

	return key
}
