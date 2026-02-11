package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/cli"
	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/core"
)

func main() {

	core.CheckForFFmpegBin()

	videoPath := cli.GetVideoPath()

	video := core.NewVideo(videoPath)

	fmt.Printf("Video metadata: %+v\n", video)
}
