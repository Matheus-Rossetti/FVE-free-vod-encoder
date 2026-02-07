package main

import (
	"fmt"

	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/cli"
	"github.com/Matheus-Rossetti/FVE-free-vod-encoder/internal/core"
)

func main() {

	core.CheckForFFmpegBin()

	videoPath := cli.GetVideoPath()

	width, height, fps := core.GetMetadata(videoPath)

	fmt.Printf("Width: %v, Height: %v, FPS: %v\n", width, height, fps)
}
