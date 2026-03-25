package dispatcher

import (
	"context"
)

type StorageProvider interface {
	Upload(context.Context, string, string) error
	HandleError([]string) error
}
