package application

import (
	"net/http"
)

func (app *App) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {

	healthCheckRes := Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": string(app.config.Env),
			"version":     Version,
		},
	}

	err := app.WriteJSON(w, http.StatusOK, healthCheckRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
