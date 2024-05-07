package application

import (
	"log/slog"
	"os"
)

func newTestingApp() *App {
	app := &App{
		Config: Config{
			Port: 4000,
			Env:  "testing",
		},
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	return app
}
