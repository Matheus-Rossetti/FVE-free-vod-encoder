package main

import (
	"os/exec"
	"runtime"
)

func (f *frevod) CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		f.slog.Error("FFmpeg not found! Frevod won't work without it :(")
		f.slog.Error("Please install FFmpeg and ensure it's in your system's PATH.")

		// TODO offer to auto install ffmpeg

		switch runtime.GOOS {
		case "windows":
			f.log.Fatal("Run the following command: winget install ffmpeg")
		case "linux":
			f.log.Println("You can install FFmpeg using your package manager")
			f.log.Println("Debian/Ubuntu: sudo apt install ffmpeg")
			f.log.Println("Fedora: sudo dnf install ffmpeg")
			f.log.Fatal("Arch: sudo pacman -S ffmpeg")
		}
	}
}
