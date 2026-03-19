package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrCreatingLocalDir     = errors.New("failed when creating local directory to store videos")
	ErrCreatingFileToCopyAt = errors.New("failed when creating a file to store output segment or manifest")
	ErrOpeningOriginalFile  = errors.New("failed when opening the original file to copy data from")
	ErrCopying              = errors.New("failed when copying data from the original file to the local copy")
)

// We call it Upload to satisfy the upload_method interface, but "Copy" would suit this better
func (l *local) Upload(ctx context.Context, key, filePath string) error {
	pathToStoreAt := filepath.Join(l.storageDir, key)

	// Create dir
	dir := filepath.Dir(pathToStoreAt)
	if err := os.MkdirAll(dir, 0755); err != nil {
		l.slog.Error(ErrCreatingLocalDir.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrCreatingLocalDir, err)
	}

	// Create file
	newFile, err := os.Create(pathToStoreAt)
	if err != nil {
		l.slog.Error(ErrCreatingFileToCopyAt.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrCreatingFileToCopyAt, err)
	}
	defer newFile.Close()

	oldFile, err := os.Open(filePath)
	if err != nil {
		l.slog.Error(ErrOpeningOriginalFile.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrOpeningOriginalFile, err)
	}
	defer oldFile.Close()

	// Copy old into new
	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		l.slog.Error(ErrCopying.Error(), "err", err)
		return fmt.Errorf("%w: %v", ErrCopying, err)
	}

	return nil
}
