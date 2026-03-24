package cli

import (
	"errors"
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
		initialErr         error
		wantErr            bool
	}{
		{
			name:               "valid URL and key",
			input:              "http://coolvideo.com key/starter",
			expectedUriType:    core.Url,
			expectedUri:        "http://coolvideo.com",
			expectedKeyStarter: "key/starter",
			initialErr:         nil,
			wantErr:            false,
		},
		{
			name:               "valid path and key",
			input:              "/coolvideo.mp4 key/starter",
			expectedUriType:    core.Path,
			expectedUri:        "/coolvideo.mp4",
			expectedKeyStarter: "key/starter",
			initialErr:         nil,
			wantErr:            false,
		},
		{
			name:               "with initial error",
			input:              "/coolvideo.mp4 key/starter",
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			initialErr:         errors.New("New Error!"),
			wantErr:            true,
		},
		{
			name:               "with uri only",
			input:              "http://coolvideo.com",
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			initialErr:         nil,
			wantErr:            true,
		},
		{
			name:               "with more than uri and key",
			input:              "http://coolvideo.com key/starter something-that-shouldn't-be-here",
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			initialErr:         nil,
			wantErr:            true,
		},
		{
			name:               "with unsupported uri",
			input:              "ftp://coolvideo.com key/starter",
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			initialErr:         nil,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cli := cli{
				slog: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				err:  testCase.initialErr,
			}

			uriType, uri, keyStarter := cli.validateInput(testCase.input)

			if uriType != testCase.expectedUriType {
				t.Errorf("")
			}

			if uri != testCase.expectedUri {
				t.Errorf("")
			}

			if keyStarter != testCase.expectedKeyStarter {
				t.Errorf("")
			}

			hasErr := (cli.err != nil)
			if hasErr != testCase.wantErr {
				t.Errorf("")
			}
		})
	}
}
