package options

import (
	"os"
	"strconv"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/go-playground/validator/v10"
	"go.yaml.in/yaml/v3"
)

func (o *options) ParseOptions() *core.Options {
	o.slog.Info("Parsing options...")
	options := &core.Options{}
	validate := validator.New()
	// TODO If config.yml file isn't found, or is malformed
	// get options from charm's Huh lib (terminal form)
	// add option to save that config into a config.yml

	if o.IsRunningInDocker() {
		o.slog.Info("Running in Docker!")

		options.Input.Cli = false
		options.Input.REST = true
		options.Upload.S3.Endpoint = os.Getenv("ENDPOINT")
		options.Upload.S3.AccessKey = os.Getenv("ACCESS_KEY")
		options.Upload.S3.SecretAccessKey = os.Getenv("SECRET_KEY")
		options.Upload.S3.BucketName = os.Getenv("BUCKET_NAME")
		if options.Upload.S3.Endpoint != "" {
			options.Upload.S3.Use = true
		}

		ssl, err := strconv.ParseBool(os.Getenv("USE_SSL"))
		if err != nil {
			o.log.Fatal("Failed to get USE_SSL")
		}
		options.Upload.S3.UseSSL = ssl

		options.Upload.Local.StoreAt = os.Getenv("STORE_AT")
		if options.Upload.Local.StoreAt != "" {
			options.Upload.Local.Use = true
		}

		concurrentEncodings, _ := strconv.Atoi(os.Getenv("CONCURRENT_ENCODINGS"))
		options.Encode.ConcurrentEncodings = concurrentEncodings

		err = validate.Struct(options)
		if err != nil {
			o.slog.Error("Validation error", "err", err)
			o.log.Fatal("Finishing program.")
		}

		return options
	}

	// Standard options
	options = &core.Options{
		Input: core.InputOptions{
			Cli:  true,
			REST: false,
		},

		Encode: core.EncodeOptions{
			ConcurrentEncodings: 2,
			OutputFFmpegCommand: false,
		},

		Upload: core.UploadOptions{
			S3: core.AmazonS3Options{
				Use:             false,
				Endpoint:        "",
				AccessKey:       "",
				SecretAccessKey: "",
				BucketName:      "",
				UseSSL:          false,
			},
			Local: core.LocalOption{
				Use:     true,
				StoreAt: "segmented-videos",
			},
		},
	}

	// TODO get filepath from flag --config

	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		o.slog.Warn("config.yml file not found, proceeding standard options.")
		return options
	}

	err = yaml.Unmarshal(configFile, options)
	if err != nil {
		o.slog.Error("Failed to parse config.yaml, proceeding with standard options.")
		return options
	}

	err = validate.Struct(options)
	if err != nil {
		o.slog.Error("Validation error", "err", err)
		o.log.Fatal("Finishing program.")
	}

	o.slog.Info("Found config.yaml, proceeding with custom options.")

	// TODO output options.

	return options
}
