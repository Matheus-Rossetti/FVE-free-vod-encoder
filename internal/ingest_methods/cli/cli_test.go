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
