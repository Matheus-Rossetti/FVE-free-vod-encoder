package core

import (
	"path/filepath"
	"strings"
)

type Video struct {

	// Source is the absolute path of where the source file (the video) is located.
	Source string
	Name   string

	Height      int
	Width       int
	AspectRatio string
	Duration    string // Seconds
	HasAudio    bool   // Useful when building the ffmpeg command | doesn't include audio tags if hasAudio == false
}

func NewVideo(videoPath string) *Video {

	metadata := GetMetadataFrom(videoPath)

	videoFileName := filepath.Base(videoPath)
	videoName := strings.TrimSuffix(videoFileName, filepath.Ext(videoFileName))

	video := Video{
		Name:        videoName,
		Source:      videoPath,
		Height:      metadata.Streams[0].Height,
		Width:       metadata.Streams[0].Width,
		AspectRatio: metadata.Streams[0].AspectRatio,
		Duration:    metadata.Streams[0].Duration,
	}

	// TODO fail if a file has more than 1 video stream | multiple video streams not supported yet
	// TODO maybe encode the first video stream and return a warning about file having more than 1 v-stream

	// check if there's at least one audio stream
	for _, stream := range metadata.Streams {
		if stream.CodecType == "audio" {
			video.HasAudio = true
			break
		} else {
			video.HasAudio = false
		}
	}

	return &video
}
