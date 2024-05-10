package application

import (
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/data"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HealthCheckPath = "/v1/model-registry/healthcheck"
	ModelRegistry   = "/v1/model-registry/"
)

const (
	Version = "1.0.0"
)

type App struct {
	config config.Config
	logger *slog.Logger
	models data.Models
}

func NewApp(cfg config.Config, logger *slog.Logger) *App {
	app := &App{
		config: cfg,
		logger: logger,
	}
	return app
}

func (app *App) Routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, HealthCheckPath, app.HealthcheckHandler)
	router.Handler(http.MethodGet, ModelRegistry, app.KubernetesClient(http.HandlerFunc(app.ModelRegistryHandler)))

	return app.RecoverPanic(app.enableCORS(router))

}
