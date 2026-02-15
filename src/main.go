package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func main() {

	core.CheckForFFmpegBin()

	// videoPath := cli.GetVideoPath()

	// video := core.NewVideo(videoPath)

	// fmt.Printf("Video metadata: %+v\n", video)

	core.Greet()
	for {
		var videoPath string

		fmt.Scanln(&videoPath)

		fmt.Printf("Initiating encoding process for %v\n", videoPath)

		go func(path string) {
			base := filepath.Base(path)
			videoName := strings.TrimSuffix(base, filepath.Ext(base))
			outputDir := filepath.Join("output", videoName)
			os.MkdirAll(outputDir, 0755)

			absoluteVideoPath, err := filepath.Abs(path)
			if err != nil {
				log.Fatal("Error getting the absolute file path for video\n", err)
			}

			cmd := core.BuildFFmpegCommand(absoluteVideoPath, videoName)
			cmd.Dir = outputDir
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal("error running the command\n", err, "for:", string(output))
			}

			fmt.Printf("%v encoded succesfully!\n", videoPath)
		}(videoPath)
	}
}
