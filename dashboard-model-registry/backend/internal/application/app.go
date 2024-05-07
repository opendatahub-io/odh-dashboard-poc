package application

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HealthCheckPath = "/v1/healthcheck"
)

const (
	Version = "1.0.0"
)

type Config struct {
	Port int
	Env  string
}

type App struct {
	Config Config
	Logger *slog.Logger
}

func NewApp(cfg Config, logger *slog.Logger) *App {
	app := &App{
		Config: cfg,
		Logger: logger,
	}
	return app
}

func (app *App) Routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, HealthCheckPath, app.HealthcheckHandler)

	return app.RecoverPanic(router)

}
