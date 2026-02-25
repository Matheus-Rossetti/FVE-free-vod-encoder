package encoder

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(workerId int, encodeQueue <-chan core.EncodeJob, filePool chan<- *os.File, options *core.Options) {
	for encodeJob := range encodeQueue {
		fmt.Printf("Encode worker %v received %v\n", workerId, encodeJob.AbsoluteVideoPath)

		video := NewVideo(encodeJob.AbsoluteVideoPath)
		cmd := BuildFFmpegCommand(video, options)

		outputDir := createOutputDir()
		// FFmpeg runs as a low-priority process, it will use 100% CPU but won't freeze the system
		cmd.Dir = outputDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal("error running the command\n", err, "for:", string(output))
		}

		filePool <- encodeJob.File // return the file to the pool

		fmt.Printf("Job %v concluded!\n", video.Name)
	}
}
