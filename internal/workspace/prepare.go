package workspace

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

func createDownloadFile() *os.File {
	os.Mkdir("downloaded-videos", 0700)

	file, err := os.CreateTemp("downloaded-videos", "download-*")
	if err != nil {
		log.Fatal("Failed when creating temp file to store video from download")
	}

	return file
}

func Prepare() (*os.File, string) {
	file := createDownloadFile()
	outputDir := createOutputDir()

	return file, outputDir
}
