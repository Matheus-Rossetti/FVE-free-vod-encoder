package core

import (
	"fmt"
	"log"
	"log/slog"
)

type frevod struct {
	Log     *log.Logger
	Slog    *slog.Logger
	Options *Options
}

func StartFrevod(log *log.Logger, slog *slog.Logger) *frevod {
	return &frevod{
		Log:  log,
		Slog: slog,
	}
}

func (f *frevod) SetOptions(options *Options) {

	options.Input.REST = true
	options.Encode.ConcurrentEncodings = 2
	options.Upload.Local.Use = true
	options.Upload.Local.StoreAt = `C:\Users\mathe\OneDrive\Desktop\segmented videos`

	f.Options = options

	f.Slog.Info(fmt.Sprintf("Using options %+v", options))
}
