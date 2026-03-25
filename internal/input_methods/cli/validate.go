package cli

import (
	"fmt"
	"strings"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (c *cli) validateInput(input string) (core.URIType, string, string) {
	if c.err != nil {
		return core.Unsupported, "", ""
	}
	if input == "" {
		c.err = fmt.Errorf("empty input")
		return core.Unsupported, "", ""
	}

	parts := strings.Fields(input)
	uri := parts[0]

	uriType, uriScheme := core.CategorizeUri(uri)
	if uriType == core.Unsupported {
		c.slog.Warn("This URI type is unsupported", "supports", "http, https, file and local paths", "got", uriScheme)
		c.err = fmt.Errorf("invalid input")
		return uriType, uri, ""
	}

	if len(parts) != 2 {
		c.slog.Warn("invalid input format", "expected", "URI KEY", "got", input)
		c.err = fmt.Errorf("invalid input")
		return core.Unsupported, "", ""
	}
	keyStarter := parts[1]

	return uriType, uri, keyStarter
}
