package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// UserWorkspaceAPIResponse represents a workspace with the user's role in it
type UserWorkspaceAPIResponse struct {
	Role      string               `json:"role"`
	Workspace WorkspaceAPIResponse `json:"workspace"`
}

// WorkspaceAPIResponse represents the response from the API for workspace-related endpoints
type WorkspaceAPIResponse struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Projects    int    `json:"projects"`
	Status      string `json:"status"`
	Users       int    `json:"users"`
}

// ListWorkspaces retrieves a list of workspaces for the current user
func (c *FCApiClient) ListWorkspaces() ([]UserWorkspaceAPIResponse, error) {
	// Make API request
	url := fmt.Sprintf("%sv1/protected/workspaces", c.host)
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
	var workspaces []UserWorkspaceAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return workspaces, nil
}

// CreateWorkspace creates a new workspace
func (c *FCApiClient) CreateWorkspace(name string, description string) (*WorkspaceAPIResponse, error) {
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
	url := fmt.Sprintf("%sv1/protected/workspaces", c.host)
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
	var workspace WorkspaceAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &workspace, nil
}
