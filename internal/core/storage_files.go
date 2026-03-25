package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func (f *frevod) createDownloadFile(dir string, index int) *os.File {
	slotName := fmt.Sprintf("%v/video-%v", dir, index)
	file, err := os.Create(slotName)
	if err != nil {
		f.Log.Fatal("Failed when creating temp file to store video from download stream", "err", err)
	}

	return file
}

func (f *frevod) PrepareStorageFiles(options *Options) []*os.File {
	dir := "storage_files"
	os.Mkdir(dir, 0700)
	var downloadSlots []*os.File
	for index := range options.Encode.ConcurrentEncodings * 2 {
		file := f.createDownloadFile(dir, index)
		downloadSlots = append(downloadSlots, file)
	}

	return downloadSlots
}

func (f *frevod) FillFilePool(files []*os.File, filePool chan<- *os.File) {
	for _, file := range files { // fill pool
		filePool <- file // TODO take this out of downloader and fill the files inside the Prepare function
	}
}

func (f *frevod) CloseAndDeleteStorageFiles(filePool <-chan *os.File) {
	for file := range filePool {
		file.Close()
		absoluteFilePath, _ := filepath.Abs(file.Name())
		os.Remove(absoluteFilePath)
		f.Slog.Warn("Closed and Deleted", "file", file.Name())
	}
}
