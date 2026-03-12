package app

import (
	"fmt"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (a *app) createDownloadFile(dir string, index int) *os.File {
	slotName := fmt.Sprintf("%v/video-%v", dir, index)
	file, err := os.Create(slotName)
	if err != nil {
		a.log.Fatal("Failed when creating temp file to store video from download stream", "err", err)
	}

	return file
}

func (a *app) PrepareStorageFiles(options *core.Options) []*os.File {
	dir := "storage_files"
	os.Mkdir(dir, 0700)
	var downloadSlots []*os.File
	for index := range options.Encode.ConcurrentEncodings * 2 {
		file := a.createDownloadFile(dir, index)
		downloadSlots = append(downloadSlots, file)
	}

	return downloadSlots
}

func (a *app) FillFilePool(files []*os.File, filePool chan<- *os.File) {
	for _, file := range files { // fill pool
		filePool <- file // TODO take this out of downloader and fill the files inside the Prepare function
	}
}
