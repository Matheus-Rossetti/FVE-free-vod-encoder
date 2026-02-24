package cli

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/downloader"
)

func Start(downloadQueue chan<- downloader.DownloadJob) {
	fmt.Println("Starting terminal input method...")
	Run(downloadQueue)
}
