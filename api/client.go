package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fusioncatalyst/paw/contracts"

	"gopkg.in/yaml.v3"
)

type FCApiClient struct {
	host          string
	authorization string
	httpClient    *http.Client
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Body)
}

// loadSettings reads and parses the fcsettings.yaml file
func loadSettings() (*contracts.SettingYAMLFile, error) {
	data, err := os.ReadFile("fcsettings.yaml")
	if err != nil {
		return nil, err
	}

	var settings contracts.SettingYAMLFile
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, errors.New("invalid settings file format: " + err.Error())
	}

	return &settings, nil
}

func NewFCApiClient() (*FCApiClient, error) {
	envHost := os.Getenv("FC_HOST")

	// Check settings file
	var fileHost string
	var hasFileHost bool

	if _, err := os.Stat("fcsettings.yaml"); err == nil {
		settings, err := loadSettings()
		if err != nil {
			return nil, errors.New("failed to load settings file: " + err.Error())
		}
		if settings.Server != "" {
			fileHost = settings.Server
			hasFileHost = true
		}
	}

	// Handle configuration conflicts and priorities
	switch {
	case envHost != "" && hasFileHost:
		return nil, errors.New("host is specified in both environment variable and settings file - please use only one source")
	case envHost != "":
		return &FCApiClient{
			host:       envHost,
			httpClient: &http.Client{},
		}, nil
	case hasFileHost:
		return &FCApiClient{
			host:       fileHost,
			httpClient: &http.Client{},
		}, nil
	default:
		return nil, errors.New("host is not specified in either environment variable (FC_HOST) or settings file")
	}
}

func (c *FCApiClient) GetHost() string {
	return c.host
}

func (c *FCApiClient) GetAuthorization() string {
	// If authorization is not set in the structure, check environment
	if c.authorization == "" {
		if envToken := os.Getenv("FC_ACCESS_TOKEN"); envToken != "" {
			// Store the token from environment in the structure
			c.authorization = envToken
		}
	}
	return c.authorization
}

// Update the authorization field setter
func (c *FCApiClient) setAuthorization(authHeader string) {
	// Remove "Bearer " prefix and any extra whitespace
	c.authorization = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
}
