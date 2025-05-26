package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type UserInfoAPIResponse struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Status string `json:"status"`
}

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
		return errors.New("failed to marshal request: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/public/users", c.host)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to read response body: " + err.Error())
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
		return errors.New("no authorization header in response")
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
		return errors.New("failed to marshal request: " + err.Error())
	}

	url := fmt.Sprintf("%sv1/public/authentication", c.host)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to read response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	authHeader := resp.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("no authorization header in response")
	}

	c.setAuthorization(authHeader)
	return nil
}

func (c *FCApiClient) GetPersonalInfo() (*UserInfoAPIResponse, error) {
	url := fmt.Sprintf("%sv1/protected/me", c.host)
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

	var userInfo UserInfoAPIResponse
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return &userInfo, nil
}
