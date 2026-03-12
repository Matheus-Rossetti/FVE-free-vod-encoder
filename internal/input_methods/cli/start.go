package cli

func (c *cli) Start() {
	c.slog.Info("Starting...")

	go c.ListenForInput()

	<-c.ctx.Done()
	c.shutdownCli()
}
