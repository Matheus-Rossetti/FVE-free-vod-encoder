package dispatcher

import (
	"context"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type StorageProvider interface {
	Name() string

	Dispatch(context.Context, core.DispatchJob) error
	HandleError(core.DispatchJob) error
}
