package options

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func ParseOptions() *core.Options {
	fmt.Println("Parsing options...")

	// Standard options
	options := &core.Options{
		Input: core.InputOptions{
			UseTerminal: false,
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
				Use:      true,
				S3UseSSL: true,
			},
			StoreLocal: "",
		},
	}

	if IsRunningInDocker() {
		options.Upload.S3.S3Endpoint = os.Getenv("S3_ENDPOINT")
		options.Upload.S3.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
		options.Upload.S3.S3AccessKey = os.Getenv("S3_SECRET_KEY")
		return options
	}

	// parse config.yaml

	log.Println("config.yml file not found, using standard options.")
	return options

}
