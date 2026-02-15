package core

import (
	"flag"
)

func GetVideoPath() string {
	videoPath := flag.String("video", "bunda", "The path of the video to be encoded")
	flag.Parse()
	return *videoPath
}
