package application

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/utils"
)

func TestHealthcheckHandler(t *testing.T) {

	app := newTestingApp()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, HealthCheckPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	app.HealthcheckHandler(rr, r)
	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal("Failed to read response body")
	}

	var healthCheckRes Envelope
	err = json.Unmarshal(body, &healthCheckRes)
	if err != nil {
		t.Fatalf("Error unmarshalling response JSON: %v", err)
	}

	expected := Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "testing",
			"version":     "1.0.0",
		},
	}

	assert.Equal(t, expected["status"], healthCheckRes["status"])
}
