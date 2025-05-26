package api

import (
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
		return "", fmt.Errorf("invalid language: %s. Must be one of: typescript, python, java, go", language)
	}

	// Make API request
	url := fmt.Sprintf("%sv1/protected/apps/%s/code/%s", c.host, appID, language)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GetAuthorization()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
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
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(code), nil
}
