package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// MessageAPIResponse represents the response from the API for message-related endpoints
type MessageAPIResponse struct {
	Description   string `json:"description"`
	Name          string `json:"name"`
	SchemaID      string `json:"schema_id"`
	SchemaVersion int    `json:"schema_version"`
}

// ListMessages retrieves a list of messages for a specific project
func (c *FCApiClient) ListMessages(projectID string) ([]MessageAPIResponse, error) {
	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/messages", c.host, projectID)
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
	var messages []MessageAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return messages, nil
}

// CreateMessage creates a new message in the specified project
func (c *FCApiClient) CreateMessage(projectID string, name string, description string, schemaID string, schemaVersion int64) (*MessageAPIResponse, error) {
	// Prepare request body
	reqBody := struct {
		Name          string `json:"name"`
		Description   string `json:"description,omitempty"`
		SchemaID      string `json:"schema_id"`
		SchemaVersion int64  `json:"schema_version"`
	}{
		Name:          name,
		Description:   description,
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/messages", c.host, projectID)
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
	var message MessageAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &message, nil
}
