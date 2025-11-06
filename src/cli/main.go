package main

import (
	"github.com/Matheus-Rossetti/video-converter-service/internal/app"
	"github.com/Matheus-Rossetti/video-converter-service/internal/video"
)

func main() {
	
	app := app.New()

	newVideo := video.NewVideo(app)
	video.Convert(newVideo, app)
}

