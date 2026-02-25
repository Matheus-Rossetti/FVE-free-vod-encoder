package encoder

import (
	"encoding/json"
	"log"
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

func GetMetadataFrom(videoPath string) *Metadata {

	// flags here could be an array, but they are few enough, so no need.
	cmd := exec.Command(
		"ffprobe",
		"-i", videoPath, // input
		"-of", "json", // output format
		"-loglevel", "error",
		"-show_entries",
		"stream=codec_type,width,height,display_aspect_ratio,duration", // used to build Video obj.
	)

	jsonOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Coudn't run ffprobe's command", err, string(jsonOutput))
	}

	var metadata Metadata
	err = json.Unmarshal(jsonOutput, &metadata)
	if err != nil {
		log.Fatal("Error parsing ffprobe's output to ffprobeOutput struct", err)
	}

	return &metadata
}
