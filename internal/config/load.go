package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/go-playground/validator/v10"
	"go.yaml.in/yaml/v4"
)

func Load(options *core.Options, path string) error {

	if isRunningInDocker() {
		fmt.Printf("\nRunning in Docker!")
		fmt.Printf("\nLoading env vars...")

		loadEnvVars(options)

		fmt.Printf("\nEnv vars loaded!\n")
		return nil
	}

	err := loadYamlConfig(options, path)
	if err != nil {
		fmt.Printf("\nCouldn't find config.yaml, using standard options...")
		loadStdOptions(options)
	}

	// TODO improve validation logs
	validate := validator.New()
	err = validate.Struct(options)
	if err != nil {
		log.Fatal("Invalid config\n", err)
	}

	return nil
}

// This is shit code, I know, I just wanna get this done already
// I'll come back to it later... or so I hope...
func loadEnvVars(options *core.Options) {
	// ingest --------
	options.Ingest.REST.Enabled = true
	options.Ingest.REST.Port = ":1137"
	// encode --------
	options.Encode.ConcurrentEncodings, _ = strconv.Atoi(os.Getenv("CONCURRENT_ENCODINGS"))
	// dispatch ------
	options.Dispatch.S3.Enabled = true
	options.Dispatch.S3.Bucket = os.Getenv("BUCKET")
	options.Dispatch.S3.Endpoint = os.Getenv("ENDPOINT")
	options.Dispatch.S3.AccessKey = os.Getenv("ACCESS_KEY")
	options.Dispatch.S3.SecretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
	options.Dispatch.S3.UseSSL, _ = strconv.ParseBool(os.Getenv("USE_SSL"))
}

func loadYamlConfig(options *core.Options, path string) error {
	if path == "" {
		path = "config.yaml"
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	fmt.Printf("\nFound config file!")
	fmt.Printf("\nLoading config from yaml...")

	err = yaml.Unmarshal(configFile, options)
	if err != nil {
		return fmt.Errorf("failed to parse yaml in %q: %w", path, err)
	}

	fmt.Printf("\nConfig loaded!\n")
	return nil
}

func loadStdOptions(options *core.Options) error {
	// ingest --------
	options.Ingest.Cli.Enabled = true
	// encode --------
	options.Encode.ConcurrentEncodings = 2
	// dispatch ------
	options.Dispatch.Local.Enabled = true
	options.Dispatch.Local.StoreAt = "."
	return nil
}
