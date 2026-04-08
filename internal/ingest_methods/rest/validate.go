package rest

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func validateBody(body io.ReadCloser) (core.URIType, string, string, error) {

	var request struct {
		Uri        string `json:"uri"`
		KeyStarter string `json:"key_starter"`
	}

	err := json.NewDecoder(body).Decode(&request)
	if err != nil {
		return core.Unsupported,
			"",
			"",
			fmt.Errorf("invalid json")
	}

	if request.Uri == "" || request.KeyStarter == "" {
		return core.Unsupported,
			"",
			"",
			fmt.Errorf("needs uri and key_starter")
	}

	uriType, scheme := core.CategorizeUri(request.Uri)
	if uriType == core.Unsupported {
		return uriType,
			"",
			"",
			fmt.Errorf("unsupported uri scheme (%v), supported schemes include: http, https, file and local paths.", scheme)
	}

	return uriType, request.Uri, request.KeyStarter, nil
}
