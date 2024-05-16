package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
)

const registerModelPath = "/registered_models"

func FetchAllRegisteredModels(client *integrations.HTTPClient) (*openapi.RegisteredModelList, error) {

	responseData, err := client.GET(registerModelPath)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered models: %w", err)
	}

	var modelList openapi.RegisteredModelList
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &modelList, nil
}

func CreateRegisteredModel(client *integrations.HTTPClient, jsonData []byte) (*openapi.RegisteredModel, error) {
	responseData, err := client.POST(registerModelPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("error posting registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}
