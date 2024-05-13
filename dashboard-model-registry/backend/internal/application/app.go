package application

import (
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/data"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HealthCheckPath      = "/api/v1/healthcheck/"
	ModelRegistry        = "/api/v1/model-registry/"
	RegisteredModelsPath = "/api/v1/model-registry/:model_registry_id/registered_models"
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

	// HTTP client routes
	router.GET(HealthCheckPath, app.HealthcheckHandler)
	router.GET(RegisteredModelsPath, app.AttachRESTClient(app.RegisteredModelsHandler))

	// Kubernetes client routes
	router.GET(ModelRegistry, app.KubernetesClient(app.ModelRegistryHandler))

	return app.RecoverPanic(app.enableCORS(router))
}
