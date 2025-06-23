package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fusioncatalyst/paw/contracts"
)

func (c *FCApiClient) ListServerResources(serverID string) ([]contracts.ResourceResponse, error) {
	url := fmt.Sprintf("%sv1/protected/servers/%s/resources", c.host, serverID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var resources []contracts.ResourceResponse
	if err := json.Unmarshal(bodyBytes, &resources); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return resources, nil
}

func (c *FCApiClient) CreateResource(serverID string, resource contracts.CreateResourceRequest) (*contracts.ResourceResponse, error) {
	url := fmt.Sprintf("%sv1/protected/servers/%s/resources", c.host, serverID)

	resource.ServerID = serverID

	jsonData, err := json.Marshal(resource)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var newResource contracts.ResourceResponse
	if err := json.Unmarshal(bodyBytes, &newResource); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &newResource, nil
}