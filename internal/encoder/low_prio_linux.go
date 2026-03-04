package encoder

import (
	"context"
	"os/exec"
)

func BuildForLowPrioExecution(ctx context.Context, ffmpegArgs []string) *exec.Cmd {
	// ------ SET LOW PRIO PROCESS FOR UNIX BASED ------

	// Use nice to set low prio
	args := append([]string{"-n", "10", "ffmpeg"}, ffmpegArgs...)
	cmd := exec.CommandContext(ctx, "nice", args...)

	return cmd
}
