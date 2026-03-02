package encoder

import (
	"os/exec"
)

func BuildForLowPrioExecution(ffmpegArgs []string) *exec.Cmd {
	// ------ SET LOW PRIO PROCESS FOR UNIX BASED ------

	// Use nice to set low prio
	args := append([]string{"-n", "10", "ffmpeg"}, ffmpegArgs...)
	cmd := exec.Command("nice", args...)

	return cmd
}
