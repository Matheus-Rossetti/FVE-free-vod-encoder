package cli

import (
	"context"
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job) {
	fmt.Println("Starting terminal input method...")
	Run(downloadQueue)
}
