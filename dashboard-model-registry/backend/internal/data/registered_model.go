package data

import (
	"encoding/json"
	"fmt"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
)

type MetadataValue struct {
	//ederign double check metatada type
	MetadataType string `json:"metadataType"`
	StringValue  string `json:"string_value"`
}

type RegisteredModel struct {
	CustomProperties map[string]MetadataValue `json:"customProperties"`
	Description      string                   `json:"description"`
	ExternalID       string                   `json:"externalId"`
	Name             string                   `json:"name"`
	ID               string                   `json:"id"`
	CreateTime       string                   `json:"createTimeSinceEpoch"`
	LastUpdateTime   string                   `json:"lastUpdateTimeSinceEpoch"`
	Owner            string                   `json:"owner"`
	State            string                   `json:"state"`
}

type RegisteredModelList struct {
	Items         []RegisteredModel `json:"items"`
	NextPageToken string            `json:"nextPageToken"`
	PageSize      int               `json:"pageSize"`
	Size          int               `json:"size"`
}

const registerModelPath = "/registered_models"

func FetchAllRegisteredModels(client *integrations.HTTPClient) (*RegisteredModelList, error) {

	responseData, err := client.GET(registerModelPath)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered models: %w", err)
	}

	var modelList RegisteredModelList
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &modelList, nil
}
