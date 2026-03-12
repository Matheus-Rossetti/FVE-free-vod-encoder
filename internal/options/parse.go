package options

import (
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (o *options) ParseOptions() *core.Options {
	o.slog.Info("Parsing options...")

	// TODO If config.yml file isn't found, or is malformed
	// get options from charm's Huh lib (terminal form)
	// add option to save that config into a config.yml

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
				S3Endpoint:        "",
				S3AccessKey:       "",
				S3SecretAccessKey: "",
				S3BucketName:      "videos",
				Use:               false,
				S3UseSSL:          true,
			},
			Local: core.LocalOption{
				Use:     true,
				StoreAt: "/home/matheusrossetti/segmented-videos",
			},
			ConcurrentUploads: 10,
		},
	}

	if o.IsRunningInDocker() {
		options.Upload.S3.S3Endpoint = os.Getenv("S3_ENDPOINT")
		options.Upload.S3.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
		options.Upload.S3.S3AccessKey = os.Getenv("S3_SECRET_KEY")
		return options
	}

	// parse config.yaml

	o.slog.Warn("config.yml file not found, using standard options.")
	return options

}
