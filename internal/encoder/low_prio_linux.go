package encoder

import (
	"os/exec"
)

func (e *encoder) BuildForLowPrioExecution(ffmpegArgs []string) *exec.Cmd {
	// ------ SET LOW PRIO PROCESS FOR UNIX BASED ------

	// Use nice to set low prio
	args := append([]string{"-n", "10", "ffmpeg"}, ffmpegArgs...)
	cmd := exec.CommandContext(e.ctx, "nice", args...)

	return cmd
}
