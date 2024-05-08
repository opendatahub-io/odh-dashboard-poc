package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	ModelRegistryNamespace = "odh-model-registries"
	ModelRegistryResource  = "modelregistries"
)

type KubernetesClient struct {
	client dynamic.Interface
}

func NewKubernetesClient() (*KubernetesClient, error) {
	// Simulate receiving an auth cookie
	cookie := "PUT-COOKIE-HERE"
	config := &rest.Config{
		BearerToken: cookie,
		Host:        "https://api.modelserving-ui.dev.datahub.redhat.com:6443",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %w", err)
	}

	return &KubernetesClient{client: dynamicClient}, nil
}

func (k *KubernetesClient) ListResources(resourceType string) ([]unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "modelregistry.opendatahub.io",
		Version:  "v1alpha1",
		Resource: resourceType,
	}

	unstructuredList, err := k.client.Resource(gvr).Namespace(ModelRegistryNamespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list custom resources in namespace %s: %v", ModelRegistryNamespace, err)
	}

	return unstructuredList.Items, nil
}
