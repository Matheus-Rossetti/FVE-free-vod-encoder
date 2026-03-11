package cli

import (
	"context"
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func Start(ctx context.Context, downloadQueue chan<- *core.Job) {
	fmt.Printf("\nStarting terminal input method...")
	go Run(downloadQueue)

	<-ctx.Done()
	shutdownCli()
}
