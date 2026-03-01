package workspace

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func createDownloadFile(index int) *os.File {
	slotName := fmt.Sprintf("download-slots/video-%v", index)
	file, err := os.Create(slotName)
	if err != nil {
		log.Fatal("Failed when creating temp file to store video from download", err)
	}

	return file
}

func Prepare(options *core.Options) []*os.File {
	os.Mkdir("download-slots", 0700)
	var downloadSlots []*os.File
	for index := range options.Encode.ConcurrentEncodings * 2 {
		file := createDownloadFile(index)
		downloadSlots = append(downloadSlots, file)
	}

	return downloadSlots
}
