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

	HasAudio            bool // Useful when building the ffmpeg command | doesn't include audio tags if hasAudio == false
	Orientation         string
	ReferenceResolution int          // vertical = width | horizontal = height | in other worlds, the size of the smaller side
	RenditionsToMake    []Resolution // if a video is 1440p/QHD, this would be -> 1440p, 1080p, 720p, 480p
}

type Resolution struct {
	name  string
	value int
}

func NewVideo(videoPath string) *Video {
	// this constructor is getting ugly

	// --------- BASE VIDEO OBJ ---------
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

	// --------- CHECK IF VIDEO HAS AUDIO ---------
	for _, stream := range metadata.Streams {
		if stream.CodecType == "audio" {
			video.HasAudio = true
			break
		} else {
			video.HasAudio = false
		}
	}

	// --------- CHECK VIDEO ORIENTATION ---------
	if video.Width > video.Height {
		video.Orientation = "horizontal"
		video.ReferenceResolution = video.Height
	} else {
		video.Orientation = "vertical" // vertical also covers square videos.
		video.ReferenceResolution = video.Width
	}

	// --------- DECIDE WHICH RENDITIONS TO MAKE ---------
	allResolutions := []Resolution{
		{"2160", 2160},
		{"1440", 1440},
		{"1080", 1080},
		{"720", 720},
		{"480", 480},
	}

	// only adds renditions of native res and lower, never upscale.
	for _, res := range allResolutions {
		if video.ReferenceResolution >= res.value {
			video.RenditionsToMake = append(video.RenditionsToMake, res)
		}
	}

	// just in case a video is smaller than 480p
	if len(video.RenditionsToMake) == 0 {
		video.RenditionsToMake = append(video.RenditionsToMake, allResolutions[len(allResolutions)-1])
	}

	return &video
}
