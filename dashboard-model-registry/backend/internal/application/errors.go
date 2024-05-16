package application

import (
	"encoding/json"
	"fmt"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
	"net/http"
	"strconv"
)

func (app *App) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	httpError := &integrations.HTTPError{
		StatusCode: http.StatusBadRequest,
		ErrorResponse: integrations.ErrorResponse{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: err.Error(),
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) errorResponse(w http.ResponseWriter, r *http.Request, error *integrations.HTTPError) {

	env := Envelope{"error": error}

	err := app.WriteJSON(w, error.StatusCode, env, nil)

	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(error.StatusCode)
	}
}

func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	httpError := &integrations.HTTPError{
		StatusCode: http.StatusInternalServerError,
		ErrorResponse: integrations.ErrorResponse{
			Code:    strconv.Itoa(http.StatusInternalServerError),
			Message: "the server encountered a problem and could not process your request",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &integrations.HTTPError{
		StatusCode: http.StatusNotFound,
		ErrorResponse: integrations.ErrorResponse{
			Code:    strconv.Itoa(http.StatusNotFound),
			Message: "the requested resource could not be found",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &integrations.HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		ErrorResponse: integrations.ErrorResponse{
			Code:    strconv.Itoa(http.StatusMethodNotAllowed),
			Message: fmt.Sprintf("the %s method is not supported for this resource", r.Method),
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {

	message, err := json.Marshal(errors)
	if err != nil {
		message = []byte("{}")
	}
	httpError := &integrations.HTTPError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorResponse: integrations.ErrorResponse{
			Code:    strconv.Itoa(http.StatusUnprocessableEntity),
			Message: string(message),
		},
	}
	app.errorResponse(w, r, httpError)
}
