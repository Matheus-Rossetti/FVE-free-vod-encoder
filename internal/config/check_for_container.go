package config

import "os"

// TODO Refactor to also check for podman
func isRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err != nil {
		return false
	}

	return true
}
