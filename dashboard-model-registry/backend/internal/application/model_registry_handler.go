package application

import (
	"net/http"

	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/models"
)

func (app *App) ModelRegistryHandler(w http.ResponseWriter, r *http.Request) {

	registries, err := models.FetchAllModelRegistry()
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
