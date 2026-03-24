package options

import "os"

// TODO Refactor to also check for podman
func (o *options) IsRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err != nil {
		return false
	}

	return true
}
