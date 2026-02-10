package core

import (
	"fmt"
	"log"
	"os/exec"
)

func GetMetadataFrom(videoPath string) string {

	input := fmt.Sprintf("-i %v", videoPath)

	// flags here could be an array, but there are few enougth, so no need
	cmd := exec.Command("ffprobe", input, "-of json", "-loglevel error", "-show_streams")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Coudn't run ffprobe's command", err)
	}

	return string(output)
}
