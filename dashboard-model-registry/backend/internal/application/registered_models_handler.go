package application

import (
	"github.com/julienschmidt/httprouter"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/data"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"net/http"
)

func (app *App) RegisteredModelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//error
	client, ok := r.Context().Value(httpClientKey).(*integrations.HTTPClient)
	if !ok {
		//fix http errors
		http.Error(w, "REST client not found", http.StatusInternalServerError)
		return
	}

	modelList, err := data.FetchAllRegisteredModels(client)
	if err != nil {
		app.serverErrorResponse(w, r, err) // Method for handling server-side errors generically
		return
	}

	modelRegistryRes := Envelope{
		"registered_models": modelList,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
