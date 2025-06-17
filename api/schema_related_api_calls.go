package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// SchemaAPIResponse represents the response from the API for schema-related endpoints
type SchemaAPIResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListSchemas retrieves a list of schemas for a specific project
func (c *FCApiClient) ListSchemas(projectID string) ([]SchemaAPIResponse, error) {
	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/schemas", c.host, projectID)
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
	var schemas []SchemaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemas); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return schemas, nil
}

// CreateSchema creates a new schema in the specified project
func (c *FCApiClient) CreateSchema(projectID string, name string, description string, schemaType string, schemaContent string) (*SchemaAPIResponse, error) {
	// First, validate that the schema content is valid JSON
	var schemaJSON interface{}
	if err := json.Unmarshal([]byte(schemaContent), &schemaJSON); err != nil {
		return nil, errors.New("invalid schema content: " + err.Error())
	}

	// Re-marshal the schema to ensure it's properly escaped
	escapedSchema, err := json.Marshal(schemaJSON)
	if err != nil {
		return nil, errors.New("failed to escape schema content: " + err.Error())
	}

	// Prepare request body
	reqBody := struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Type        string `json:"type"`
		Schema      string `json:"schema"`
	}{
		Name:        name,
		Description: description,
		Type:        schemaType,
		Schema:      string(escapedSchema),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/projects/%s/schemas", c.host, projectID)
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
	var schema SchemaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&schema); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &schema, nil
}

// UpdateSchema updates an existing schema
func (c *FCApiClient) UpdateSchema(schemaID string, schemaContent string) (*SchemaAPIResponse, error) {
	// First, validate that the schema content is valid JSON
	var schemaJSON interface{}
	if err := json.Unmarshal([]byte(schemaContent), &schemaJSON); err != nil {
		return nil, errors.New("invalid schema content: " + err.Error())
	}

	// Re-marshal the schema to ensure it's properly escaped
	escapedSchema, err := json.Marshal(schemaJSON)
	if err != nil {
		return nil, errors.New("failed to escape schema content: " + err.Error())
	}

	// Prepare request body
	reqBody := struct {
		Schema string `json:"schema"`
	}{
		Schema: string(escapedSchema),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal request: " + err.Error())
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/schemas/%s", c.host, schemaID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
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
	var schema SchemaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&schema); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &schema, nil
}
