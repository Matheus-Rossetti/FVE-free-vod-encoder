package main

import (
	"os/exec"
	"runtime"
)

func (f *frevod) CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		f.slog.Error("FFmpeg binary not found! Frevod wont work without it :(")
		f.slog.Error("Please install FFmpeg and ensure it's in your system's PATH.")

		switch runtime.GOOS {
		case "windows":
			f.log.Fatal("Here's the download link: https://www.ffmpeg.org/download.html")
		case "linux":
			f.log.Println("You can install FFmpeg using your package manager")
			f.log.Println("Debian/Ubuntu: sudo apt install ffmpeg")
			f.log.Println("Fedora: sudo dnf install ffmpeg")
			f.log.Fatal("Arch: sudo pacman -S ffmpeg")
		}
	}
}
