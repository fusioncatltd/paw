package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ProjectAPIResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *FCApiClient) ListProjects() ([]ProjectAPIResponse, error) {
	url := fmt.Sprintf("%sv1/protected/projects", c.host)
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

	var projects []ProjectAPIResponse
	if err := json.Unmarshal(bodyBytes, &projects); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return projects, nil
}

func (c *FCApiClient) CreateProject(
	name, description string,
	createdByType string,
	createdByID string,
	isPrivate bool,
) (*ProjectAPIResponse, error) {
	reqBody := struct {
		Name          string `json:"name"`
		Description   string `json:"description,omitempty"`
		CreatedByType string `json:"created_by_type"`
		CreatedByID   string `json:"created_by_id"`
		IsPrivate     bool   `json:"is_private"`
	}{
		Name:          name,
		Description:   description,
		CreatedByType: createdByType,
		CreatedByID:   createdByID,
		IsPrivate:     isPrivate,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/protected/projects", c.host)
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

	var project ProjectAPIResponse
	if err := json.Unmarshal(bodyBytes, &project); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &project, nil
}

// ImportProject uploads a file to the specified project and processes the import
func (c *FCApiClient) ImportProject(projectID string, filePath string) error {
	// First, verify and read the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("failed to read file: " + err.Error())
	}

	// Create request body with file content as text
	reqBody := struct {
		YAML string `json:"yaml"`
	}{
		YAML: string(fileContent),
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("failed to marshal request: " + err.Error())
	}

	// Create the request
	url := fmt.Sprintf("%sv1/protected/projects/%s/imports", c.host, projectID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to read response body: " + err.Error())
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	return nil
}

// GenerateCode generates code for a specific application in a project
func (c *FCApiClient) GenerateCode(projectID string, appID string) (*ProjectAPIResponse, error) {
	url := fmt.Sprintf("%sv1/protected/projects/%s/apps/%s/generate", c.host, projectID, appID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
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

	var project ProjectAPIResponse
	if err := json.Unmarshal(bodyBytes, &project); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &project, nil
}
