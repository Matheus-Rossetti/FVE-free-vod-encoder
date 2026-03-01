package options

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"

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

	// try to get options from config.yml
	found, config := GetOptionsFromFile()
	if found {
		log.Println("Found config file.")
		yaml.Unmarshal(config, options)
		return options
	}

	// get config from env vars

	log.Println("Config file not found, using standard options.")
	return options
}
