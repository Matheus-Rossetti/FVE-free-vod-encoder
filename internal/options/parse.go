package options

import (
	"fmt"
	"log"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func ParseOptions() *core.Options {
	fmt.Println("Parsing options...")

	// Standard options
	options := &core.Options{
		Input: core.InputOptions{
			UseTerminal: true,
			UseREST:     true,
		},

		Encode: core.EncodeOptions{
			Codec:               "h.264",
			SegmentType:         "fmp4",
			SegmentDuration:     2,
			ConcurrentEncodings: 2,
			OutputFFmpegCommand: false,
		},

		Upload: core.UploadOptions{
			S3: core.AmazonS3Options{
				Use: false,
			},
			StoreLocal: "",
		},
	}

	if IsRunningInDocker() {
		return options
	}

	log.Println("config.yml file not found, using standard options.")
	return options

}
