package cli

import (
	"fmt"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (c *cli) validateInput(input string) (core.URIType, string, string, error) {
	if input == "" {
		return core.Unsupported, "", "", fmt.Errorf("empty input")
	}

	// Validates URI first
	parts := strings.Fields(input)
	uri := parts[0]

	uriType, uriScheme := core.CategorizeUri(uri)
	if uriType == core.Unsupported {
		c.slog.Warn("This URI type is unsupported", "supports", "http, https, file and local paths", "got", uriScheme)
		return uriType, uri, "", fmt.Errorf("unsupported uri")
	}

	// Then key
	if len(parts) != 2 {
		c.slog.Warn("invalid input format", "expected", "URI KEY", "got", input)
		return uriType, uri, "", fmt.Errorf("input needs two parts")
	}
	keyStarter := parts[1]

	return uriType, uri, keyStarter, nil
}
