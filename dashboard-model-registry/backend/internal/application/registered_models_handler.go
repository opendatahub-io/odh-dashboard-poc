package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/data"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/validation"
	"net/http"
)

func (app *App) GetRegisteredModelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	client, ok := r.Context().Value(httpClientKey).(*integrations.HTTPClient)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	modelList, err := data.FetchAllRegisteredModels(client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

func (app *App) CreateRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(*integrations.HTTPClient)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var model openapi.RegisteredModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	if err := validation.ValidateRegisteredModel(model); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(model)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
		return
	}

	createdModel, err := data.CreateRegisteredModel(client, jsonData)
	if err != nil {
		var httpErr *integrations.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if createdModel == nil {
		app.serverErrorResponse(w, r, fmt.Errorf("created model is nil"))
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", RegisteredModelsPath, *createdModel.Id))
	app.WriteJSON(w, http.StatusCreated, createdModel, nil)
}
