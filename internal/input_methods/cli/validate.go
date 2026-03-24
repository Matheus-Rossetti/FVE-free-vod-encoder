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

	parts := strings.Fields(input)
	if len(parts) != 2 {
		c.slog.Warn("invalid input format", "expected", "URI KEY", "got", input)
		c.err = fmt.Errorf("invalid input")
		return core.Unsupported, "", ""
	}

	uri := parts[0]
	keyStarter := parts[1]

	uriType, uriScheme := core.CategorizeUri(uri)
	if uriType == core.Unsupported {
		c.slog.Warn("This URI type is unsupported", "supports", "http, https and file", "got", uriScheme)
		c.err = fmt.Errorf("invalid input")
		return core.Unsupported, "", ""
	}

	return uriType, uri, keyStarter
}
