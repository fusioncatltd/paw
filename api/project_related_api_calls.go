package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var projects []ProjectAPIResponse
	if err := json.Unmarshal(bodyBytes, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
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
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%sv1/protected/projects", c.host)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var project ProjectAPIResponse
	if err := json.Unmarshal(bodyBytes, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &project, nil
}
