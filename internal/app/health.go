package app

import (
	"os/exec"
	"runtime"
)

func (a *app) CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		a.slog.Error("FFmpeg binary not found! Frevod wont work without it :(")
		a.slog.Error("Please install FFmpeg and ensure it's in your system's PATH.")

		switch runtime.GOOS {
		case "windows":
			a.log.Fatal("Here's the download link: https://www.ffmpeg.org/download.html")
		case "linux":
			a.log.Println("You can install FFmpeg using your package manager")
			a.log.Println("Debian/Ubuntu: sudo apt install ffmpeg")
			a.log.Println("Fedora: sudo dnf install ffmpeg")
			a.log.Fatal("Arch: sudo pacman -S ffmpeg")
		}
	}
}
