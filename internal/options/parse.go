package options

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

// TODO set standard values and substitute them for the ones in .env
func ParseOptions() *core.Options {

	fmt.Println("Parsing options...")

	return &core.Options{
		Codec:                "h.264",
		SegmentType:          "fmp4",
		SegmentDuration:      2,
		UseTerminal:          false,
		UseREST:              true,
		UseRabbitMQ:          false,
		StoreLocal:           true,
		StoreS3:              false,
		Upload:               false,
		ExposeMetrics:        false,
		ConcurrentEncodings:  2, // also specifies the amount of downloaded videos waiting to be processed
		OutputFFmpegCommand:  false,
		OutputEncodedVideoTo: "output",
		MaxStoredVideos:      2,
	}
}
