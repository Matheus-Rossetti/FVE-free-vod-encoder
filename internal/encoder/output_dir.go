package encoder

import (
	"log"
	"os"
)

func createOutputDir() string {
	os.Mkdir("output", 0700)

	// The return value here is "output/video-{random-string}"
	dirName, err := os.MkdirTemp("output", "video-*")
	if err != nil {
		log.Fatal("Error creating output directory: ", err)
	}

	return dirName
}
