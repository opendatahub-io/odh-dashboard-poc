package application

import (
	"github.com/julienschmidt/httprouter"
	k8s "github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"net/http"
)

func (app *App) ModelRegistryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	client, ok := r.Context().Value(k8sClientKey).(*k8s.KubernetesClient)
	if !ok {
		//ederign fix this to the right http error
		http.Error(w, "Kubernetes client not found", http.StatusInternalServerError)
		return
	}

	registries, err := app.models.ModelRegistry.FetchAllModelRegistry(client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := Envelope{
		"model_registry": registries,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
