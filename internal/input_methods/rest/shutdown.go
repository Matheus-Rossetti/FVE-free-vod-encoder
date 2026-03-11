package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func shutdownServer(server *http.Server) {
	fmt.Printf("\nREST input method shutting down...")

	serverCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := server.Shutdown(serverCtx)
	if err != nil {
		fmt.Printf("\nError shutting down REST input method: %v", err)
		return
	}
}
