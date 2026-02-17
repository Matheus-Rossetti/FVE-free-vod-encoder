package core

import (
	"log"
	"os/exec"
	"runtime"
)

func CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("FFmpeg binary not found! Frevod wont work without it :(")
		log.Printf("Please install FFmpeg and ensure it's in your system's PATH.")

		switch runtime.GOOS {
		case "windows":
			log.Fatal("Here's the download link: https://www.ffmpeg.org/download.html")
		case "linux":
			log.Println("You can install FFmpeg using your package manager")
			log.Println("Debian/Ubuntu: sudo apt install ffmpeg")
			log.Println("Fedora: sudo dnf install ffmpeg")
			log.Fatal("Arch: sudo pacman -S ffmpeg")
		}
	}

	log.Printf("FFmpeg binary found! Frevod is ready to go :)")
}
