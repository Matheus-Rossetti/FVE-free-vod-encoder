package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/core"
	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/io/cli"
)

func main() {

	core.CheckForFFmpegBin()

	// videoPath := cli.GetVideoPath()

	// video := core.NewVideo(videoPath)

	// fmt.Printf("Video metadata: %+v\n", video)

	cli.Greet()
	for {
		var videoPath string

		fmt.Scanln(&videoPath)

		fmt.Printf("Initiating encoding process for %v\n", videoPath)

		go func(path string) {
			base := filepath.Base(path)
			name := strings.TrimSuffix(base, filepath.Ext(base))
			outputDir := filepath.Join("HLS", name)
			os.MkdirAll(outputDir, 0755)

			cmd := core.BuildFFmpegCommand(path, outputDir)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("error running the command\n%v \nfor: %v\n", err, string(output))
			}

			fmt.Printf("%v encoded succesfully!\n", videoPath)
		}(videoPath)
	}
}
