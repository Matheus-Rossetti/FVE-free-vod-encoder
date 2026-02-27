package uploader

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func Start() {

	// will need base s3 path to upload to, maybe get from input method
	// eg: channel-123/crazy video number 6/
	// then after it I'll put streams_x and the normal and master playlilsts

	filepath.WalkDir("output/video-3294969873", printAllFiles)

}

// print file
func printAllFiles(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.IsDir() {
		return nil
	}

	relativePath, err := filepath.Rel("output/video-3294969873", path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(relativePath)
	filename := filepath.Base((relativePath))

	fmt.Printf("%v ||| %v\n", dir, filename)

	return nil
}
