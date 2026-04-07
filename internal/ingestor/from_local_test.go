package ingestor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAbsolutePath(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "With prefix",
			path:    "file://my/cool/path",
			wantErr: false,
		},
		{
			name:    "Without prefix",
			path:    "/my/cool/path",
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			_, err := getAbsolutePath(tt.path)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
			}
		})
	}
}

func TestCheckExistence(t *testing.T) {

	tempDir := t.TempDir()
	file, _ := os.Create(filepath.Join(tempDir, "file_that_exists"))
	file.Close()

	testCases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Exists",
			path:    filepath.Join(tempDir, "file_that_exists"),
			wantErr: false,
		},
		{
			name:    "Ghost file",
			path:    filepath.Join(tempDir, "maybe_in_another_universe"),
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			_, err := checkExistence(tt.path)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
			}
		})
	}
}

// Do you know how to unit test this? Please help.
func TestCheckFileForVideo(t *testing.T) {

}
