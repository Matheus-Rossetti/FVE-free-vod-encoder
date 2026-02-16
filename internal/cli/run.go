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

	var videoPath string
	for {
		_, err := fmt.Scanln(&videoPath)
		if err != nil {
			log.Fatal("Error getting input from the terminal", err)
		}
		// TODO check if URI is url or path

		// if it's url
		// TODO download

		id := "1"

		absoluteVideoPath, err := filepath.Abs(videoPath)
		if err != nil {
			log.Fatal("Couldn't get absolute path for video: ", videoPath, "\nError: ", err)
		}

		jobQueue <- core.VideoJob{
			Source:    "Terminal",
			VideoPath: absoluteVideoPath, // absolute path is safer than relative path
			Id:        id,
		}
	}
}
