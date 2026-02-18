package core

import (
	"fmt"
	"log"
	"os"
)

func CreateOutputDir(videoName string) string {

	// Make basic output dir
	os.Mkdir("output", 0700)

	// concatenate video.nome with "-*", the -> * <- is where MkdirTemp adds a random string
	dirPattern := fmt.Sprintf("%v-*", videoName)

	// create the dir inside output/ The return value here is "output/{video.Name}-{random-string}"
	dirName, err := os.MkdirTemp("output", dirPattern)
	if err != nil {
		log.Fatal("Error creating output directory: ", err)
	}

	return dirName
}
