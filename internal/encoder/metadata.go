package encoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type Metadata struct {
	Streams []struct {
		CodecType   string `json:"codec_type"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
		AspectRatio string `json:"display_aspect_ratio"`
		Duration    string `json:"duration"`
	} `json:"streams"`
}

var (
	ErrRunningFFprobe       = errors.New("failed when running FFprobe")
	ErrParsingFFprobeOutput = errors.New("failed parsing FFprobe's output to Metadata struct")
)

func GetMetadataFrom(videoPath string) (*Metadata, error) {

	// flags here could be an array, but they are few enough, so no need.
	cmd := exec.Command(
		"ffprobe",
		"-i", videoPath,
		"-of", "json", // output format
		"-loglevel", "error",
		"-show_entries",
		"stream=codec_type,width,height,display_aspect_ratio,duration", // what we use to build Video obj.
	)

	jsonOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRunningFFprobe, string(jsonOutput))
	}

	var metadata Metadata
	err = json.Unmarshal(jsonOutput, &metadata)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingFFprobeOutput, err)
	}

	return &metadata, nil
}
