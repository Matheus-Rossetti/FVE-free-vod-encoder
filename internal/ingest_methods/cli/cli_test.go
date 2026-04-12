package cli

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func TestValidateInput(t *testing.T) {
	testCases := []struct {
		name               string
		input              string
		expectedUriType    core.URIType
		expectedUri        string
		expectedKeyStarter string
		wantErr            bool
	}{
		{
			name:               "twoTokens_http_ok",
			input:              "http://coolvideo.com key/starter",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "key/starter",
			wantErr:            false,
		},
		{
			name:               "twoTokens_path_ok",
			input:              "/coolvideo.mp4 key/starter",
			expectedUriType:    core.Path,
			expectedUri:        "/coolvideo.mp4",
			expectedKeyStarter: "key/starter",
			wantErr:            false,
		},
		{
			name:               "oneToken_err",
			input:              "http://coolvideo.com",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "threeTokens_err",
			input:              "http://coolvideo.com key/starter something-that-shouldn't-be-here",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "twoTokens_ftp_err",
			input:              "ftp://coolvideo.com key/starter",
			expectedUriType:    core.Unsupported,
			expectedUri:        "ftp://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "empty_err",
			input:              "",
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cli := cli{
				slog: slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}

			uriType, uri, keyStarter, err := cli.validateInput(testCase.input)

			if testCase.expectedUriType != uriType {
				t.Errorf("expected URIType %v, got %v", testCase.expectedUriType, uriType)
			}

			if testCase.expectedUri != uri {
				t.Errorf("expected uri %v, got %v", testCase.expectedUri, uri)
			}

			if testCase.expectedKeyStarter != keyStarter {
				t.Errorf("expected keystarter %v, got %v", testCase.expectedKeyStarter, keyStarter)
			}

			hasErr := (err != nil)
			if testCase.wantErr != hasErr {
				t.Errorf("expected error to be %v, got %v", testCase.wantErr, hasErr)
			}
		})
	}
}

func TestStart_Integration(t *testing.T) {
	t.Run("twoTokens_http_ok", func(t *testing.T) {
		opts := core.NewOptions()
		opts.Encode.ConcurrentEncodings = 1
		ingest, _, _ := core.StartQueues(opts)

		oldStdin := os.Stdin
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdin = r
		t.Cleanup(func() {
			os.Stdin = oldStdin
			_ = r.Close()
			_ = w.Close()
		})

		ctx, cancel := context.WithCancel(context.Background())
		discardLog := log.New(io.Discard, "", 0)
		discardSlog := slog.New(slog.NewTextHandler(io.Discard, nil))
		c := NewCli(ctx, discardLog, discardSlog)

		done := make(chan struct{})
		go func() {
			c.Start()
			close(done)
		}()

		if _, err := w.WriteString("http://coolvideo.com key/starter\n"); err != nil {
			t.Fatal(err)
		}

		select {
		case job := <-ingest:
			if job.DownloadJob.UriType != core.Url {
				t.Errorf("UriType: got %v, want Url", job.DownloadJob.UriType)
			}
			if job.DownloadJob.VideoUri != "http://coolvideo.com" {
				t.Errorf("VideoUri: got %q, want http://coolvideo.com", job.DownloadJob.VideoUri)
			}
			if job.UploadJob.KeyStarter != "key/starter" {
				t.Errorf("KeyStarter: got %q, want key/starter", job.UploadJob.KeyStarter)
			}
			if job.DownloadJob.Source != "cli" {
				t.Errorf("Source: got %q, want cli", job.DownloadJob.Source)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("timed out waiting for job on ingest queue")
		}

		cancel()

		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("Start did not return after context cancel")
		}
	})
}
