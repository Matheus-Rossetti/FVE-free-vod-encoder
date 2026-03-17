package encoder

import (
	"os/exec"
	"runtime"
	"syscall"
)

func (e *encoder) BuildForLowPrioExecution(ffmpegArgs []string) *exec.Cmd {
	// ------ SET LOW PRIO PROCESS FOR WINDOWS ------
	cmd := exec.CommandContext(e.ctx, "ffmpeg", ffmpegArgs...)

	if runtime.GOOS == "windows" {
		const MICROSOFT_MAGIC_CONSTANT_THAT_STARTS_LOW_PRIORITY_PROCESSES = 0x00004000

		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: MICROSOFT_MAGIC_CONSTANT_THAT_STARTS_LOW_PRIORITY_PROCESSES,
		}
	}

	return cmd
}
