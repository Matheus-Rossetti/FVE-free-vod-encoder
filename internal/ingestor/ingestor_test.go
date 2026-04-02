package ingestor

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPrepareFileForDownload(t *testing.T) {
	testCases := []struct {
		name      string
		fileState string
		wantErr   bool
	}{
		{
			name:      "opened file",
			fileState: "open",
			wantErr:   false,
		},
		{
			name:      "closed file",
			fileState: "closed",
			wantErr:   true,
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

			if tt.fileState == "closed" {
				file.Close()
			}

			ingestor := ingestor{}

			err = ingestor.prepareFileForDownload(file)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
				return // no need to continue test
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

func TestCheckFileType(t *testing.T) {

	// mock server
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

			}))
	defer server.Close()

	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "video file type",
			url:     server.URL,
			wantErr: false,
		},
		{
			name:    "not video file type",
			url:     server.URL,
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

		})
	}

}
