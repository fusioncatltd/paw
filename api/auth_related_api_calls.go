package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *FCApiClient) SignUp(email, password string) error {
	type SignupRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody := SignupRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%sv1/public/users", c.host)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	// Get and store the authorization header
	authHeader := resp.Header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("no authorization header in response")
	}

	// Use the new setter method instead of direct assignment
	c.setAuthorization(authHeader)
	return nil
}

func (c *FCApiClient) SignIn(email, password string) error {
	type SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody := SignInRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%sv1/public/authentication", c.host)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	authHeader := resp.Header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("no authorization header in response")
	}

	c.setAuthorization(authHeader)
	return nil
}
