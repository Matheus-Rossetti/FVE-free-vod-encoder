package cli

import (
	"fmt"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/Matheus-Rossetti/frevod/internal/workspace"
)

func Run(downloadQueue chan<- core.Job) {
	fmt.Println("Frevod is waiting for paths or URLs in the terminal!")
	fmt.Println("Just paste it down below and the video will be encoded!")

	var videoUri string
	for {
		fmt.Print("\n>> ")
		_, err := fmt.Scanln(&videoUri)
		if err != nil {
			log.Fatal("Error getting input from the terminal", err)
		}

		file, outputDir := workspace.Prepare()
		job := core.NewJob(file, outputDir, videoUri, "terminal")
		downloadQueue <- *job
	}
}
