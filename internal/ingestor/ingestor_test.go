package ingestor

import (
	"os"
	"testing"
)

func TestPrepareFileForDownload(t *testing.T) {
	testCases := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "sucessfully prepare file",
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "temp-file-*")
			if err != nil {
				t.Errorf("failed to create a temp file to run tests with")
			}
			defer file.Close()
			defer os.Remove(file.Name())

			err = prepareFileForDownload(file)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
			}

			meta, _ := file.Stat()
			if meta.Size() != 0 {
				t.Errorf("expected file size to be 0, got %v", meta.Size())
			}

			offSet, _ := file.Seek(0, 1)
			if offSet != 0 {
				t.Errorf("expected offset of file to be at 0, got %v", offSet)
			}
		})
	}
}

// I couldn't write this test, so I'm just gonna leave it in god's hands
func TestCheckBodyForVideo(t *testing.T) {
}
