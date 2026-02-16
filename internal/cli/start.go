package cli

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(jobQueue chan<- core.VideoJob) {
	fmt.Println("Starting terminal input method...")

	Run(jobQueue)
}
