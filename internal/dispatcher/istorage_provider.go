package dispatcher

import (
	"context"
)

type StorageProvider interface {
	Dispatch(context.Context, string, string) error
	HandleError([]string) error
}
