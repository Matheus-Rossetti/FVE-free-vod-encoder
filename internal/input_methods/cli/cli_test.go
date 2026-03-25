package cli

import (
	"log/slog"
	"os"
	"testing"

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
			name:               "valid url and key",
			input:              "http://coolvideo.com key/starter",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "key/starter",
			wantErr:            false,
		},
		{
			name:               "valid path and key",
			input:              "/coolvideo.mp4 key/starter",
			expectedUriType:    core.Path,
			expectedUri:        "/coolvideo.mp4",
			expectedKeyStarter: "key/starter",
			wantErr:            false,
		},
		{
			name:               "with just uri (no key)",
			input:              "http://coolvideo.com",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "with more than url and key",
			input:              "http://coolvideo.com key/starter something-that-shouldn't-be-here",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "with unsupported uri",
			input:              "ftp://coolvideo.com key/starter",
			expectedUriType:    core.Unsupported,
			expectedUri:        "ftp://coolvideo.com",
			expectedKeyStarter: "",
			wantErr:            true,
		},
		{
			name:               "with empty input",
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

func TestPushJob(t *testing.T) {
	testCases := []struct {
		name        string
		source      string
		uriType     core.URIType
		uri         string
		keyStarter  string
		expectedJob core.Job
	}{
		{
			name:       "correct job from cli",
			uriType:    core.Url,
			uri:        "https://mycoolvideo.com",
			keyStarter: "key/starter",
			expectedJob: core.Job{
				DownloadJob: core.DownloadJob{
					Source:   "cli",
					UriType:  core.Url,
					VideoUri: "https://mycoolvideo.com",
				},
				UploadJob: core.UploadJob{
					KeyStarter: "key/starter",
				},
			},
		},
	}

	for _, tt := range testCases {

		queue := make(chan *core.Job, 2)

		cli := cli{
			downloadQueue: queue,
		}

		cli.pushJob(tt.uriType, tt.uri, tt.keyStarter)

		job := <-queue

		if tt.expectedJob.DownloadJob.Source != job.DownloadJob.Source {
			t.Errorf("expected source to be %v, got %v",
				tt.expectedJob.DownloadJob.Source,
				job.DownloadJob.Source)
		}

		if tt.expectedJob.DownloadJob.UriType != job.DownloadJob.UriType {
			t.Errorf("expected uri type to be %v, got %v",
				tt.expectedJob.DownloadJob.UriType,
				job.DownloadJob.UriType)
		}

		if tt.expectedJob.DownloadJob.VideoUri != job.DownloadJob.VideoUri {
			t.Errorf("expected uri to be %v, got %v",
				tt.expectedJob.DownloadJob.VideoUri,
				job.DownloadJob.VideoUri)
		}

		if tt.expectedJob.UploadJob.KeyStarter != job.UploadJob.KeyStarter {
			t.Errorf("expected key starter to be %v, got %v",
				tt.expectedJob.UploadJob.KeyStarter,
				job.UploadJob.KeyStarter)
		}
	}
}
