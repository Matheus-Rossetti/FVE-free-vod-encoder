package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// We call it Upload to satisfy the upload_method interface, but "Copy" would suit this better
func (l *Local) Upload(ctx context.Context, key, filePath string) error {
	pathToStoreAt := filepath.Join(l.storageDir, key)

	// Create dir
	dir := filepath.Dir(pathToStoreAt)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("\nError creating directories: %v", err)
		return err
	}

	// Create file
	newFile, err := os.Create(pathToStoreAt)
	if err != nil {
		fmt.Printf("\nError creating file to copy segment to: %v", err)
		return err
	}
	defer newFile.Close()

	oldFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("\nError opening origin file to copy segment from: %v", err)
		return err
	}
	defer oldFile.Close()

	// Copy old into new
	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		fmt.Printf("\nError copying file: %v", err)
		return err
	}

	return nil
}
