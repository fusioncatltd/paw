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

func (c *FCApiClient) CreateServer(server *contracts.CreateServerRequest) (*contracts.Server, error) {
	jsonData, err := json.Marshal(server)
	if err != nil {
		return nil, errors.New("failed to marshal server data: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/protected/servers", c.host)
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
	url := fmt.Sprintf("%sv1/protected/servers", c.host)
	if projectID != "" {
		url = fmt.Sprintf("%s?project_id=%s", url, projectID)
	}

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

	var result contracts.ServersListResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &result, nil
}

func (c *FCApiClient) GetServer(serverID string) (*contracts.Server, error) {
	url := fmt.Sprintf("%sv1/protected/servers/%s", c.host, serverID)
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

	var result contracts.Server
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &result, nil
}

func (c *FCApiClient) UpdateServer(serverID string, update *contracts.UpdateServerRequest) (*contracts.Server, error) {
	jsonData, err := json.Marshal(update)
	if err != nil {
		return nil, errors.New("failed to marshal update data: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/protected/servers/%s", c.host, serverID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK {
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

func (c *FCApiClient) DeleteServer(serverID string) error {
	url := fmt.Sprintf("%sv1/protected/servers/%s", c.host, serverID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("failed to execute request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	return nil
}