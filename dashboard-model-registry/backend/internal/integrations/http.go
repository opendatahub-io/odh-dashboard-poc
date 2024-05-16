package integrations

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
	token   string
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HTTPError struct {
	StatusCode int `json:"-"`
	ErrorResponse
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s - %s", e.StatusCode, e.Code, e.Message)
}

func NewHTTPClient(env config.Environment, modelRegistryBaseURL string) (*HTTPClient, error) {
	var token string

	switch env {
	case config.Development:
		cmd := exec.Command("oc", "whoami", "--show-token")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %s, error: %w", output, err)
		}
		token = strings.TrimSpace(string(output))
	case config.Production:
		/* Production: Extract token from request headers
		       //   ederign TODO
		   		//token = r.Header.Get("Authorization")
		   		//if token == "" {
		   		//	return nil, fmt.Errorf("no Authorization token provided")
		   		//}

		   		//token = strings.TrimPrefix(token, "Bearer ")


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
	default:
		return nil, fmt.Errorf("unsupported environment: %v", env)
	}

	return &HTTPClient{
		client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		baseURL: modelRegistryBaseURL,
		token:   token,
	}, nil
}

func (c *HTTPClient) GET(url string) ([]byte, error) {
	fullURL := c.baseURL + url
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}

func (c *HTTPClient) POST(url string, body io.Reader) ([]byte, error) {
	fullURL := c.baseURL + url
	fmt.Println(fullURL)
	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token)

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if response.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling error response: %w", err)
		}
		httpError := &HTTPError{
			StatusCode:    response.StatusCode,
			ErrorResponse: errorResponse,
		}
		//Sometimes the code comes empty from model registry API
		//also not all error codes are correctly implemented
		//see https://github.com/kubeflow/model-registry/issues/95
		if httpError.ErrorResponse.Code == "" {
			httpError.ErrorResponse.Code = strconv.Itoa(response.StatusCode)
		}
		return nil, httpError
	}

	return responseBody, nil
}
