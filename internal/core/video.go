package core

import (
	"encoding/json"
	"log"
)

type Stream struct {
	CodecType string `json:"codec_type"`
}

type Video struct {
	Streams     []Stream `json:"streams"` //
	Height      int      `json:"height"`
	Width       int      `json:"width"`
	AspectRatio string   `json:"aspect_ratio"`
	Duration    float32  `json:"duration"` //seconds
	HasAudio    bool     // Useful when building the ffmpeg command | doesn't include audio tags if hasAudio == false
}

func NewVideo(videoPath string) *Video {

	var video Video
	jsonMetadata := GetMetadataFrom(videoPath)

	// map metadata to video struct
	err := json.Unmarshal([]byte(jsonMetadata), &video)
	if err != nil {
		log.Fatal("Couldn't map ffprobe output to video struct.", err)
	}

	// check if there's at least one audio stream
	for _, stream := range video.Streams {
		if stream.CodecType == "audio" {
			video.HasAudio = true
			break
		} else {
			video.HasAudio = false
		}
	}

	return &video
}
