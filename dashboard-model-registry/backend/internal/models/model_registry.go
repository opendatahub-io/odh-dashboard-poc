package models

import (
	"fmt"

	k8s "github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/integrations"
)

type ModelRegistry struct {
	Name string `json:"name"`
}

func FetchAllModelRegistry() ([]ModelRegistry, error) {
	client, err := k8s.NewKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %w", err)
	}

	resources, err := client.ListResources(k8s.ModelRegistryResource)
	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries []ModelRegistry
	for _, item := range resources {
		fmt.Println(item)
		registry := ModelRegistry{
			Name: item.GetName(),
		}
		registries = append(registries, registry)
	}

	return registries, nil
}
