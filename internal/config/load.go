package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"github.com/go-playground/validator/v10"
	"go.yaml.in/yaml/v4"
)

func Load(options *core.Options, path string) error {

	err := loadStdOptions(options)
	if err != nil {
		return err
	}

	if isRunningInDocker() {
		err := loadEnvVars(options)
		if err != nil {
			return err
		}
		return nil
	}

	err = loadYaml(options, path)
	if err != nil {
		return err
	}

	// TODO improve validation logs
	validate := validator.New()
	err = validate.Struct(options)
	if err != nil {
		log.Fatal("Unvalid config\n", err)
	}

	return nil
}

func loadStdOptions(options *core.Options) error {
	options.Encode.ConcurrentEncodings = 2
	return nil
}

func loadEnvVars(options *core.Options) error {
	return nil
}

func loadYaml(options *core.Options, path string) error {
	if path == "" {
		path = "config.yaml"
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	err = yaml.Unmarshal(configFile, options)
	if err != nil {
		return fmt.Errorf("failed to parse yaml in %q: %w", path, err)
	}

	return nil
}
