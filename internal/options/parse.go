package options

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func ParseOptions() *core.Options {
	fmt.Println("Parsing options...")

	// TODO If config.yml file isn't found, or is malformed
	// get options from charm's Huh lib (terminal form)
	// add option to save that config into a config.yml

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
				S3Endpoint:        "",
				S3AccessKey:       "",
				S3SecretAccessKey: "",
				S3BucketName:      "videos",
				Use:               false,
				S3UseSSL:          true,
			},
			StoreLocal:        "",
			ConcurrentUploads: 10,
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
