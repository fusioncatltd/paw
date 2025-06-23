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

func (c *FCApiClient) CreateServer(projectID string, server *contracts.CreateServerRequest) (*contracts.Server, error) {
	// Create request body with only the fields the API expects
	reqBody := struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description,omitempty"`
	}{
		Name:        server.Name,
		Type:        server.Type,
		Description: server.Description,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal server data: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/protected/projects/%s/servers", c.host, projectID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute request: " + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var result contracts.Server
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &result, nil
}

func (c *FCApiClient) ListServers(projectID string) (*contracts.ServersListResponse, error) {
	if projectID == "" {
		return nil, errors.New("project ID is required")
	}
	
	url := fmt.Sprintf("%sv1/protected/projects/%s/servers", c.host, projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute request: " + err.Error())
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

	var servers []contracts.Server
	if err := json.Unmarshal(bodyBytes, &servers); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	result := &contracts.ServersListResponse{
		Servers: servers,
		Total:   len(servers),
	}

	return result, nil
}


