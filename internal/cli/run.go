package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Run(jobQueue chan<- core.VideoJob) {
	fmt.Println("Frevod is waiting for paths or URLs in the terminal!")
	fmt.Println("Just paste it down below and the video will be encoded!")

	var videoUri string
	for {
		fmt.Print("\n>> ")
		_, err := fmt.Scanln(&videoUri)
		if err != nil {
			log.Fatal("Error getting input from the terminal", err)
		}

		// uriType := core.GetUriType(videoUri)
		// if uriType == "url" {
		// 	videoUri = core.DownloadAndStoreVideo(videoUri)
		// }

		videoUri = core.DownloadAndStoreVideo(videoUri)

		absoluteVideoPath, err := filepath.Abs(videoUri)
		if err != nil {
			log.Fatal("Couldn't get absolute path for video: ", videoUri, "\nError: ", err)
		}

		jobQueue <- core.VideoJob{
			Source:            "Terminal",
			AbsoluteVideoPath: absoluteVideoPath, // absolute path is safer than relative path
		}
	}
}
