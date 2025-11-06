package app

import "github.com/Matheus-Rossetti/video-converter-service/internal/flags"

type App struct {
	Flags *flags.Flags
}

func New() *App {
	flags := flags.New()
	flags.Bind()

	return &App{
		Flags: flags,
	}
}
