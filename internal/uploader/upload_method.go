package uploader

import (
	"context"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type StorageProvider interface {
	Upload(context.Context, *core.Job, string) error
	HandleError()
}
