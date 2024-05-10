package application

import (
	"context"
	"fmt"
	k8s "github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"log"
	"net/http"
)

type contextKey string

const k8sClientKey contextKey = "k8sClient"

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
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (app *App) KubernetesClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		client, err := k8s.NewKubernetesClient(app.config.Env)
		if err != nil {
			//ederign fix also this error
			log.Printf("Failed to create Kubernetes client: %v", err)
			http.Error(w, "Failed to create Kubernetes client", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), k8sClientKey, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
