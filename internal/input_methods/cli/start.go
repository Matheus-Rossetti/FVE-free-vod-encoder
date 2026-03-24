package cli

import (
	"bufio"
	"os"
	"time"
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
			uriType, uri, keyStarter := c.validateInput(input)
			c.pushJob(uri, uriType, keyStarter)

			if c.err != nil {
				c.slog.Warn("ignoring input", "why", c.err.Error())
				c.err = nil
				continue
			}
		}
	}()

	<-c.ctx.Done()
	time.Sleep(time.Second * 2)
	c.slog.Info("Shutting down...")
}
