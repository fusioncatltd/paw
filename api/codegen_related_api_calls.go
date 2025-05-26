package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GenerateAppCode generates code for an application in the specified language
func (c *FCApiClient) GenerateAppCode(appID string, language string) (string, error) {
	// Validate language
	validLanguages := map[string]bool{
		"typescript": true,
		"python":     true,
		"java":       true,
		"go":         true,
	}
	if !validLanguages[language] {
		return "", errors.New("invalid language: " + language + ". Must be one of: typescript, python, java, go")
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/apps/%s/code/%s", c.host, appID, language)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	// Read response body
	code, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read response body: " + err.Error())
	}

	return string(code), nil
}
