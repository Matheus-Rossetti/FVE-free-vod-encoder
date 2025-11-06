package tests

import (
	"fmt"

	"github.com/Matheus-Rossetti/video-converter-service/internal/video"
)

func TestGetMetadata() {
	width, height, fps := video.GetMetadata("video de teste.mp4")

	fmt.Printf("Width: %v\n Height: %v\n Fps: %v\n", width, height, fps)
}
