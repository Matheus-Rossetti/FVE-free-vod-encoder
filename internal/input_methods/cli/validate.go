package cli

import (
	"fmt"
	"strings"
)

func (c *cli) validateInput(input string) (string, string) {
	if c.err != nil {
		return "", ""
	}

	parts := strings.Fields(input)

	if len(parts) != 2 {
		c.slog.Warn("invalid input format", "expected", "URI KEY", "got", input)
		c.err = fmt.Errorf("invalid input")
		return "", ""
	}

	uri := parts[0]
	keyStarter := parts[1]
	return uri, keyStarter
}
