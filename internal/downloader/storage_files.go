package downloader

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func createDownloadFile(dir string, index int) *os.File {
	slotName := fmt.Sprintf("/video-%v", index)
	file, err := os.Create(slotName)
	if err != nil {
		log.Fatal("Failed when creating temp file to store video from download stream", err)
	}

	return file
}

func PrepareStorageFiles(options *core.Options) []*os.File {
	dir := "storage_files"
	os.Mkdir(dir, 0700)
	var downloadSlots []*os.File
	for index := range options.Encode.ConcurrentEncodings * 2 {
		file := createDownloadFile(dir, index)
		downloadSlots = append(downloadSlots, file)
	}

	return downloadSlots
}
