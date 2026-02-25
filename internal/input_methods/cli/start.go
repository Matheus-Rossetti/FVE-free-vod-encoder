package cli

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(downloadQueue chan<- core.Job) {
	fmt.Println("Starting terminal input method...")
	Run(downloadQueue)
}
