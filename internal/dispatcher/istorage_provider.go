package dispatcher

import (
	"context"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

// each provider should deal with errors by themselves
type StorageProvider interface {
	Name() string
	Dispatch(context.Context, core.DispatchJob) error
}
