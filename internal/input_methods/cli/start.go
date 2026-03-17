package cli

import "time"

func (c *cli) Start() {
	c.slog.Info("Starting...")

	go c.ListenForInput()

	<-c.ctx.Done()
	time.Sleep(time.Second * 2)
	c.slog.Info("Shutting down...")
}
