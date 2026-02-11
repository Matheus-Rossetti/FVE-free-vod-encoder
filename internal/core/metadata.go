package core

import (
	"log"
	"os/exec"
)

func GetMetadataFrom(videoPath string) string {

	// flags here could be an array, but there are few enougth, so no need
	cmd := exec.Command("ffprobe", "-i", videoPath, "-of", "json", "-loglevel", "error", "-show_streams")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Coudn't run ffprobe's command", err, string(output))
	}

	return string(output)
}
