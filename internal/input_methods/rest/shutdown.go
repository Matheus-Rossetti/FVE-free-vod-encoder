package rest

import (
	"context"
	"net/http"
	"time"
)

func (r *rest) shutdownServer(server *http.Server) {
	r.slog.Info("Shutting down...")

	serverCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := server.Shutdown(serverCtx)
	if err != nil {
		r.slog.Error("Couldn't shutdown!", "err", err)
		return
	}
}
