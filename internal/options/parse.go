package options

import (
	"fmt"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func ParseOptions() *core.Options {

	fmt.Println("Parsing options...")
	options = &core.Options{}

	return &core.Options{
		Codec:                "h.264",
		SegmentType:          "fmp4",
		SegmentDuration:      2,
		UseTerminal:          true,
		UseREST:              true,
		UseRabbitMQ:          false,
		StoreLocal:           true,
		UseS3:                true,
		S3Endpoint:           "",
		S3AccessKey:          "",
		S3SecretAccessKey:    "",
		S3UseSSL:             true,
		ExposeMetrics:        false,
		ConcurrentEncodings:  2, // also specifies the amount of downloaded videos waiting to be processed
		OutputFFmpegCommand:  false,
		OutputEncodedVideoTo: "output",
		MaxStoredVideos:      2,

		ConcurrentUploads: 10,
	}
}
