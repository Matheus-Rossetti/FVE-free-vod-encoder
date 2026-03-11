package core

import (
	"os"
)

type EncodeJob struct {
	AbsoluteVideoPath string
	File              *os.File
	DownloadedFile    bool
}
