package config

import (
	"fmt"
	"os"

	"github.com/Matheus-Rossetti/frevod/internal/core"
	"go.yaml.in/yaml/v4"
)

func Load(options *core.Options, path string) error {

	err := standardValues(options)
	if err != nil {
		return err
	}

	if isRunningInDocker() {
		err := loadEnvVars(options)
		if err != nil {
			return err
		}
	}

	err = loadYaml(options, path)
	if err != nil {
		return err
	}

	return nil
}

func standardValues(target any) error {
	_ = target
	return nil
}

func loadEnvVars(target any) error {
	_ = target
	return nil
}

func loadYaml(target any, path string) error {
	if path == "" {
		path = "config.yaml"
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(configFile, target); err != nil {
		return fmt.Errorf("failed to parse yaml in %q: %w", path, err)
	}

	return nil
}
