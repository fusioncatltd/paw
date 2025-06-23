package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fusioncatalyst/paw/api"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
)

func CreateServer(ctx context.Context, cmd *cli.Command) error {
	name := cmd.String("name")
	serverType := cmd.String("type")
	description := cmd.String("description")
	projectID := cmd.String("project-id")

	if name == "" {
		return errors.New("Server name is required")
	}
	if serverType == "" {
		return errors.New("Server type is required")
	}
	if description == "" {
		return errors.New("Server description is required")
	}
	if projectID == "" {
		return errors.New("Project ID is required")
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize API client: %v", err))
	}

	req := &contracts.CreateServerRequest{
		Name:        name,
		Type:        serverType,
		Description: description,
		ProjectID:   projectID,
	}


	result, err := client.CreateServer(projectID, req)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(fmt.Sprintf("Failed to create server: %s", apiErr))
		}
		return errors.New(fmt.Sprintf("Failed to create server: %v", err))
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to format response: %v", err))
	}
	fmt.Println(string(jsonData))

	return nil
}

func ListServers(ctx context.Context, cmd *cli.Command) error {
	projectID := cmd.String("project-id")

	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize API client: %v", err))
	}

	result, err := client.ListServers(projectID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(fmt.Sprintf("Failed to list servers: %s", apiErr))
		}
		return errors.New(fmt.Sprintf("Failed to list servers: %v", err))
	}

	if len(result.Servers) == 0 {
		fmt.Println("No servers found")
		return nil
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to format response: %v", err))
	}
	fmt.Println(string(jsonData))

	return nil
}

