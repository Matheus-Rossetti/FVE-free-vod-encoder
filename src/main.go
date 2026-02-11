package main

import (
	"fmt"
	"log"

	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/core"
)

func main() {

	core.CheckForFFmpegBin()

	// videoPath := cli.GetVideoPath()

	// // video := core.NewVideo(videoPath)

	// fmt.Printf("Video metadata: %+v\n", video)

	cmd := core.BuildFFmpegCommand()

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error running the command\n", err, string(output))
	}

	fmt.Print(string(output))
}
