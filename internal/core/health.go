package core

import (
	"fmt"
	"log"
	"os/exec"
)

func CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("FFmpeg binary not found! FVE wont work without it :(")
		log.Printf("Please install FFmpeg and ensure it's in your system's PATH.")
		log.Fatal("Here's the download link: https://www.ffmpeg.org/download.html")
		fmt.Scanln()
	}

}
