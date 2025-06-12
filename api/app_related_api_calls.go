package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// AppAPIResponse represents the response from the API for app-related endpoints
type AppAPIResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListApps retrieves a list of apps for a specific project
func (c *FCApiClient) ListApps(projectID string) ([]AppAPIResponse, error) {
	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/apps", c.host, projectID)
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

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	// Read and parse response body
	var apps []AppAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return apps, nil
}

// CreateApp creates a new app in the specified project
func (c *FCApiClient) CreateApp(projectID string, name string, description string) (*AppAPIResponse, error) {
	// Prepare request body
	reqBody := struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}{
		Name:        name,
		Description: description,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/apps", c.host, projectID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	// Read and parse response body
	var app AppAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &app, nil
}
