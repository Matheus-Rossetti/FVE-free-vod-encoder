package core

import (
	"fmt"
	"log"
	"os"
)

func CreateOutputDir(videoName string, options *Options) string {

	outputDir := options.OutputEncodedVideoTo

	// Make basic output dir
	os.Mkdir(outputDir, 0700)

	// concatenate video.nome with "-*", the -> * <- is where MkdirTemp adds a random string
	dirPattern := fmt.Sprintf("%v-*", videoName)

	// create the dir inside {outputDir}/ The return value here is "{outputDir}/{video.Name}-{random-string}"
	dirName, err := os.MkdirTemp("output", dirPattern)
	if err != nil {
		log.Fatal("Error creating output directory: ", err)
	}

	return dirName
}

func CreateDownloadDir() {
	os.Mkdir("downloaded-videos", 0700)
}
