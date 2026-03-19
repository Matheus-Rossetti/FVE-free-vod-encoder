package uploader

import (
	"context"
)

type StorageProvider interface {
	Upload(context.Context, string, string) error
	HandleError([]string) error
}
