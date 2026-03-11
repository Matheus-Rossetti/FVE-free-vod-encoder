package local

import (
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

type Local struct {
	storageDir string
}

func Start(options *core.Options) *Local {

	err := os.MkdirAll(options.Upload.Local.StoreAt, 0700)
	if err != nil {
		log.Fatalf("Couldn't create dir to store videos to %v", err)
	}

	return &Local{
		storageDir: options.Upload.Local.StoreAt,
	}
}
