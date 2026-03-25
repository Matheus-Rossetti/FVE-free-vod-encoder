package core

import (
	"os/exec"
	"runtime"
)

func (f *frevod) CheckForFFmpegBin() {
	cmd := exec.Command("ffmpeg", "-version")
	_, err := cmd.CombinedOutput()

	if err != nil {
		f.Slog.Error("FFmpeg not found! Frevod won't work without it :(")
		f.Slog.Error("Please install FFmpeg and ensure it's in your system's PATH.")

		// TODO offer to auto install ffmpeg

		switch runtime.GOOS {
		case "windows":
			f.Log.Fatal("Run the following command: winget install ffmpeg")
		case "linux":
			f.Log.Println("You can install FFmpeg using your package manager")
			f.Log.Println("Debian/Ubuntu: sudo apt install ffmpeg")
			f.Log.Println("Fedora: sudo dnf install ffmpeg")
			f.Log.Fatal("Arch: sudo pacman -S ffmpeg")
		}
	}
}
