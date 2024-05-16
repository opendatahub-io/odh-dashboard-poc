package application

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	integrations "github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"net/http"
)

type contextKey string

const k8sClientKey contextKey = "k8sClient"

const httpClientKey contextKey = "httpClientKey"

func (app *App) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO ederign restrict CORS to a much smaller set of trusted origins.
		//deal with preflight requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (app *App) KubernetesClient(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		client, err := integrations.NewKubernetesClient(app.config.Env)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), k8sClientKey, client)
		next(w, r.WithContext(ctx), ps)
	}
}

func (app *App) AttachRESTClient(handler func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		modelRegistryID := ps.ByName("model_registry_id")

		modelRegistryBaseURL := resolveModelRegistryURL(modelRegistryID)
		client, err := integrations.NewHTTPClient(app.config.Env, modelRegistryBaseURL)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), httpClientKey, client)
		handler(w, r.WithContext(ctx), ps)
	}
}

// TEMP METHOD!!!!
func resolveModelRegistryURL(id string) string {
	//ederign fix this hardcoded usage
	if id == "internal-modelregistry" {
		return "http://internal-modelregistry-http-odh-model-registries.apps.modelserving-ui.dev.datahub.redhat.com/api/model_registry/v1alpha3"
	}
	return "http://modelresgistry-sample-http-odh-model-registries.apps.modelserving-ui.dev.datahub.redhat.com/api/model_registry/v1alpha3"
}
