package integrations

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
		   export const USER_ACCESS_TOKEN = 'x-forwarded-access-token';
		   		    const accessToken = request.headers[USER_ACCESS_TOKEN];
		           if (!accessToken) {
		             fastify.log.error(
		               `No ${USER_ACCESS_TOKEN} header. Cannot make a pass through call as this user.`,
		             );
		             throw new Error('No access token provided by oauth. Cannot make any API calls to kube.');
		           }
		           headers = {
		             ...kubeHeaders,
		             Authorization: `Bearer ${accessToken}`,
		           };

		*/
		return nil, errors.New(fmt.Sprintf("%s is not implemented yet", config.Production))
	default:
		return nil, errors.New("environment is not defined")
	}

	config := &rest.Config{
		BearerToken: token,
		Host:        ModelRegistryK8sHost,
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

	unstructuredList, err := k.client.Resource(gvr).Namespace(ModelRegistryK8sNamespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list custom resources in namespace %s: %v", ModelRegistryK8sNamespace, err)
	}

	return unstructuredList.Items, nil
}
