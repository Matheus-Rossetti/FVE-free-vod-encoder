package ingestor

import (
	"os"
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
	testCases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Exists",
			path:    "file_that_exists",
			wantErr: false,
		},
		{
			name:    "Ghost file",
			path:    "/maybe_in_another_universe",
			wantErr: true,
		},
	}

	file, _ := os.Create("file_that_exists")
	defer file.Close()
	defer os.Remove(file.Name())

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
