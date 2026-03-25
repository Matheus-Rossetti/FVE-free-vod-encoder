package cli

import (
	"bufio"
	"os"
	"time"
)

func (c *cli) Start() {
	go func() {
		c.slog.Info("Frevod is waiting for URIs in your terminal! Use: uri key")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				break
			}

			uriType, uri, keyStarter, err := c.validateInput(scanner.Text())
			if err != nil {
				c.slog.Error("ignoring input", "why", err.Error())
				continue
			}

			c.pushJob(uriType, uri, keyStarter)
		}
	}()

	<-c.ctx.Done()
	time.Sleep(time.Second * 2)
	c.slog.Info("Shutting down...")
}
