package core

import (
	"log"
	"os/exec"
)

func RunFFmpeg(cmd *exec.Cmd, outputDir string) {

	cmd.Dir = outputDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error running the command\n", err, "for:", string(output))
	}
}
