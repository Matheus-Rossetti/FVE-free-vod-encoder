package options

import (
	"os"
	"path/filepath"
)

func GetOptionsFromFile() (bool, []byte) {
	// look in the current directory
	info, err := os.Stat("config.yml")
	if err != nil {
		return false, nil
	}

	configFilePath, _ := filepath.Abs(info.Name())

	config, err := os.ReadFile(configFilePath)
	if err != nil {
		return false, nil
	}
	return true, config
}

func UseEnvVars() {

}
