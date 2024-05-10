package application

import (
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"log/slog"
	"os"
)

func newTestingApp() *App {
	app := &App{
		config: config.Config{
			Port: 4000,
			Env:  "testing",
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	return app
}
