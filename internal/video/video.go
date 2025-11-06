package video

import (
	"log"

	"github.com/Matheus-Rossetti/video-converter-service/internal/app"
)

type Video struct {
	Path   string
	Name   string
	Width  int
	Height int
	Fps    int
}

func NewVideo(app *app.App) *Video {
	path := app.Flags.InputPath
	name := app.Flags.OutputName
	width, height, fps := GetMetadata(path)

	if path == "" {
		log.Fatal("Missing path, use --path")
	}

	return &Video{
		Path:   path,
		Name:   name,
		Width:  width,
		Height: height,
		Fps:    fps,
	}
}
