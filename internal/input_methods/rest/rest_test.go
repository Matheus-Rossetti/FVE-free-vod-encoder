package rest

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func TestValidateBody(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		expectedUriType    core.URIType
		expectedUri        string
		expectedKeyStarter string
		expectedStatus     int
		wantErr            bool
	}{
		{
			name:               "valid body",
			body:               `{"uri": "http://mycoolvideo.com", "key_starter": "key/starter"}`,
			expectedUriType:    core.Url,
			expectedUri:        "http://mycoolvideo.com",
			expectedKeyStarter: "key/starter",
			expectedStatus:     http.StatusOK,
			wantErr:            false,
		},
		{
			name:               "empty body",
			body:               ``,
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			expectedStatus:     http.StatusUnprocessableEntity,
			wantErr:            true,
		},
		{
			name:               "lacking key_starter",
			body:               `{"uri": "http://mycoolvideo.com"}`,
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			expectedStatus:     http.StatusUnprocessableEntity,
			wantErr:            true,
		},
		{
			name:               "unsupported uri",
			body:               `{"uri": "ftp://mycoolvideo.com", "key_starter": "key/starter"}`,
			expectedUriType:    core.Unsupported,
			expectedUri:        "",
			expectedKeyStarter: "",
			expectedStatus:     http.StatusUnprocessableEntity,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {

		body := io.NopCloser(strings.NewReader(tt.body))

		uriType, uri, keyStarter, status, err := validateBody(body)

		if tt.expectedUriType != uriType {
			t.Errorf("expected uri type to be %v, got %v", tt.expectedUri, uriType)
		}

		if tt.expectedUri != uri {
			t.Errorf("expected uri to be %v, got %v", tt.expectedUri, uri)
		}

		if tt.expectedKeyStarter != keyStarter {
			t.Errorf("expected key starter to be %v, got %v", tt.expectedKeyStarter, keyStarter)
		}

		if tt.expectedStatus != status {
			t.Errorf("expected status to be %v, got %v", tt.expectedStatus, status)
		}

		hasErr := (err != nil)
		if tt.wantErr != hasErr {
			t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
		}
	}
}
