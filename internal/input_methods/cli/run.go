package cli

import (
	"fmt"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Run(downloadQueue chan<- core.DownloadJob) {
	fmt.Println("Frevod is waiting for paths or URLs in the terminal!")
	fmt.Println("Just paste it down below and the video will be encoded!")

	var videoUri string
	for {
		fmt.Print("\n>> ")
		_, err := fmt.Scanln(&videoUri)
		if err != nil {
			log.Fatal("Error getting input from the terminal", err)
		}

		downloadJob := core.NewDownloadJob(videoUri, "terminal")
		downloadQueue <- downloadJob
	}
}
