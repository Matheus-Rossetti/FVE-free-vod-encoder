package core

import (
	"encoding/json"
	"log"
)

type Stream struct {
	CodecType   string `json:"codec_type"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	AspectRatio string `json:"display_aspect_ratio"`
	Duration    string `json:"duration"`
}

type Streams struct {
	Index []Stream `json:"streams"`
}

type Video struct {
	Height      int
	Width       int
	AspectRatio string
	Duration    string //seconds
	HasAudio    bool   // Useful when building the ffmpeg command | doesn't include audio tags if hasAudio == false
}

func NewVideo(videoPath string) *Video {

	var streams Streams
	jsonMetadata := GetMetadataFrom(videoPath)

	// map metadata to video struct
	err := json.Unmarshal([]byte(jsonMetadata), &streams)
	if err != nil {
		log.Fatal("Couldn't map ffprobe output to video struct.", err)
	}

	video := Video{
		Height:      streams.Index[0].Height,
		Width:       streams.Index[0].Width,
		Duration:    streams.Index[0].Duration,
		AspectRatio: streams.Index[0].AspectRatio,
	}

	// check if there's at least one audio stream
	for _, stream := range streams.Index {
		if stream.CodecType == "audio" {
			video.HasAudio = true
			break
		} else {
			video.HasAudio = false
		}
	}

	return &video
}
