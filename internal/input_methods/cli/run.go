package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/input_methods"
)

func Run(downloadQueue chan<- *core.Job) {
	fmt.Println("Frevod is waiting for paths or URLs in the terminal!")
	fmt.Println("Just paste it down below and the video will be encoded!")

	var videoUri string
	for {
		time.Sleep(time.Second)
		fmt.Printf("\nURI >> ")
		_, err := fmt.Scanln(&videoUri)
		if err != nil {
			log.Fatal("Error getting input from the terminal", err)
		}

		uriType := input_methods.CategorizeUri(videoUri)
		if uriType == "unsupported uri" {
			fmt.Printf("\nUnsupported Uri")
			continue
		}

		job := core.NewJob()
		job.DownloadJob.Source = "Terminal"
		job.DownloadJob.UriType = uriType
		job.DownloadJob.VideoUri = videoUri

		job.UploadJob.S3KeyStarter = "videotest/"

		downloadQueue <- job
	}
}
