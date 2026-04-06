package encoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrRunningFFprobe       = errors.New("failed when running FFprobe")
	ErrParsingFFprobeOutput = errors.New("failed parsing FFprobe's output")
	ErrGettingMetadata      = errors.New("failed to get metadata from video")
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

type Metadata struct {
	Streams []struct {
		CodecType   string `json:"codec_type"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
		AspectRatio string `json:"display_aspect_ratio"`
		Duration    string `json:"duration"`
	} `json:"streams"`
}

type Resolution struct {
	name  string
	value int
}

func (e *encoder) NewVideo(videoPath string) (*Video, error) {
	// this constructor is getting ugly

	// --------- BASE VIDEO OBJ ---------
	metadata, err := e.GetMetadataFrom(videoPath)
	if err != nil {
		e.slog.Error(ErrGettingMetadata.Error(), "err", err)
		return nil, fmt.Errorf("%w: %v", ErrGettingMetadata, err)
	}

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

	// just in case a video res is lower than the lowest supported res
	if len(video.RenditionsToMake) == 0 {
		video.RenditionsToMake = append(video.RenditionsToMake, allResolutions[len(allResolutions)-1])
	}

	return &video, nil
}

func (e *encoder) GetMetadataFrom(videoPath string) (*Metadata, error) {

	// flags here could be an array, but they are few enough, so no need.
	cmd := exec.Command(
		"ffprobe",
		"-i", videoPath,
		"-of", "json", // output format
		"-loglevel", "error",
		"-show_entries",
		"stream=codec_type,width,height,display_aspect_ratio,duration", // used to build Video obj.
	)

	jsonOutput, err := cmd.CombinedOutput()
	if err != nil {
		e.slog.Error(ErrRunningFFprobe.Error(), "err", err)
		return nil, fmt.Errorf("%w: %v", ErrRunningFFprobe, string(jsonOutput))
	}

	var metadata Metadata
	err = json.Unmarshal(jsonOutput, &metadata)
	if err != nil {
		e.slog.Error(ErrParsingFFprobeOutput.Error(), "err", err)
		return nil, fmt.Errorf("%w: %v", ErrParsingFFprobeOutput, err)
	}

	return &metadata, nil
}
