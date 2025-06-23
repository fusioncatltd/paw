package actions

import (
	"context"
	"encoding/json"
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
		return cli.Exit("Server name is required", 1)
	}
	if serverType == "" {
		return cli.Exit("Server type is required", 1)
	}
	if description == "" {
		return cli.Exit("Server description is required", 1)
	}
	if projectID == "" {
		return cli.Exit("Project ID is required", 1)
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
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
			return cli.Exit(fmt.Sprintf("Failed to create server: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to create server: %v", err), 1)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format response: %v", err), 1)
	}
	fmt.Println(string(jsonData))

	return nil
}

func ListServers(ctx context.Context, cmd *cli.Command) error {
	projectID := cmd.String("project-id")

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	result, err := client.ListServers(projectID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to list servers: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to list servers: %v", err), 1)
	}

	if len(result.Servers) == 0 {
		fmt.Println("No servers found")
		return nil
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format response: %v", err), 1)
	}
	fmt.Println(string(jsonData))

	return nil
}

