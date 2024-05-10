package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"log"
	"os/exec"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	ModelRegistryNamespace = "odh-model-registries"
	ModelRegistryResource  = "modelregistries"
	ModelRegistryHost      = "https://api.modelserving-ui.dev.datahub.redhat.com:6443"
)

type KubernetesClient struct {
	client dynamic.Interface
}

func NewKubernetesClient(env config.Environment) (*KubernetesClient, error) {
	var token string
	switch env {
	case config.Development:
		log.Println("Setting up development kubernetes environment")
		cmd := exec.Command("oc", "whoami", "--show-token")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %s, error: %w", output, err)
		}
		token = strings.TrimSpace(string(output))
	case config.Staging:
		return nil, errors.New(fmt.Sprintf("%s is not implemented yet", config.Staging))
	case config.Production:
		/* ederign TODO
		   const getCurrentToken = async (currentUser: User) => {
		     return new Promise((resolve, reject) => {
		       if (currentContext === 'inClusterContext') {
		         const location =
		           currentUser?.authProvider?.config?.tokenFile ||
		           '/var/run/secrets/kubernetes.io/serviceaccount/token';
		         fs.readFile(location, 'utf8', (err, data) => {
		           if (err) {
		             reject(err);
		           }
		           resolve(data);
		         });
		       } else {
		         resolve(currentUser?.token || '');
		       }
		     });
		   };


		*/
		return nil, errors.New(fmt.Sprintf("%s is not implemented yet", config.Production))
	default:
		return nil, errors.New("environment is not defined")
	}

	config := &rest.Config{
		BearerToken: token,
		Host:        ModelRegistryHost,
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
