package ingestor

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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

			meta, _ := file.Stat()
			if meta.Size() != 0 {
				t.Errorf("expected file size to be 0, got %v", meta.Size())
			}

			offSet, _ := file.Seek(0, 1)
			if offSet != 0 {
				t.Errorf("expected offset of file to be at 0, got %v", offSet)
			}

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
			}
		})
	}
}

// I couldn't write this test, so I'm just gonna leave it in god's hands
func TestCheckBodyForVideo(t *testing.T) {
}

func TestMakeRequest(t *testing.T) {

	testCases := []struct {
		name       string
		timeout    bool
		statusCode int
		wantErr    bool
	}{
		{
			name:       "Successfuly crate and execute request",
			timeout:    false,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Status not 2**",
			timeout:    false,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "Request Timeout",
			timeout:    true,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			testServer := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						time.Sleep(time.Millisecond * 20)
						w.WriteHeader(tt.statusCode)
						fmt.Fprint(w, "hello from the test server!")
					},
				),
			)
			defer testServer.Close()

			ctx := context.Background()
			if tt.timeout {
				ctx, _ = context.WithTimeout(ctx, time.Millisecond*10)
			}

			_, err := makeRequest(testServer.URL, ctx)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
			}

		})
	}
}
