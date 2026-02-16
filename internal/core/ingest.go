package core

import (
	"os"
	"path/filepath"
)

// TODO add error handling
func CreateOutputDir(directoryName string) string {

	outputDir := filepath.Join("output", directoryName)
	os.MkdirAll(outputDir, 0755)

	return outputDir
}
