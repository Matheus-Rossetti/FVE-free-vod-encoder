package cli

import (
	"bufio"
	"os"
	"time"

	"github.com/Matheus-Rossetti/frevod/internal/core"
)

func (c *cli) Start() {
	c.slog.Info("Starting...")

	go func() {
		c.slog.Info("Frevod is waiting for URIs in your terminal! Use: videoUri key")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			uri, keyStarter := c.validateInput(input)
			uriType := core.CategorizeUri(uri)
			c.pushJob(uri, uriType, keyStarter)

			if c.err != nil {
				c.slog.Warn("ignoring input", "why", c.err.Error())
				c.err = nil
				continue
			}
			time.Sleep(time.Second) // avoid spam
		}
	}()

	<-c.ctx.Done()
	time.Sleep(time.Second * 2)
	c.slog.Info("Shutting down...")
}
