package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fusioncatalyst/paw/api"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
)

func CreateServer(ctx context.Context, cmd *cli.Command) error {
	name := cmd.String("name")
	serverType := cmd.String("type")
	description := cmd.String("description")
	projectID := cmd.String("project-id")
	resources := cmd.String("resources")
	binds := cmd.String("binds")

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

	if resources != "" {
		var serverResources []contracts.ServerResource
		if err := json.Unmarshal([]byte(resources), &serverResources); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to parse resources JSON: %v", err), 1)
		}
		req.Resources = serverResources
	}

	if binds != "" {
		var serverBinds []contracts.ServerBind
		if err := json.Unmarshal([]byte(binds), &serverBinds); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to parse binds JSON: %v", err), 1)
		}
		req.Binds = serverBinds
	}

	result, err := client.CreateServer(req)
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

func GetServer(ctx context.Context, cmd *cli.Command) error {
	serverID := cmd.String("server-id")

	if serverID == "" {
		return cli.Exit("Server ID is required", 1)
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	result, err := client.GetServer(serverID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to get server: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to get server: %v", err), 1)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format response: %v", err), 1)
	}
	fmt.Println(string(jsonData))

	return nil
}

func UpdateServer(ctx context.Context, cmd *cli.Command) error {
	serverID := cmd.String("server-id")
	name := cmd.String("name")
	serverType := cmd.String("type")
	description := cmd.String("description")
	resources := cmd.String("resources")
	binds := cmd.String("binds")

	if serverID == "" {
		return cli.Exit("Server ID is required", 1)
	}

	if name == "" && serverType == "" && description == "" && resources == "" && binds == "" {
		return cli.Exit("At least one field to update must be provided", 1)
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	req := &contracts.UpdateServerRequest{}

	if name != "" {
		req.Name = name
	}
	if serverType != "" {
		req.Type = serverType
	}
	if description != "" {
		req.Description = description
	}

	if resources != "" {
		var serverResources []contracts.ServerResource
		if err := json.Unmarshal([]byte(resources), &serverResources); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to parse resources JSON: %v", err), 1)
		}
		req.Resources = serverResources
	}

	if binds != "" {
		var serverBinds []contracts.ServerBind
		if err := json.Unmarshal([]byte(binds), &serverBinds); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to parse binds JSON: %v", err), 1)
		}
		req.Binds = serverBinds
	}

	result, err := client.UpdateServer(serverID, req)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to update server: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to update server: %v", err), 1)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format response: %v", err), 1)
	}
	fmt.Println(string(jsonData))

	return nil
}

func DeleteServer(ctx context.Context, cmd *cli.Command) error {
	serverID := cmd.String("server-id")
	force := cmd.Bool("force")

	if serverID == "" {
		return cli.Exit("Server ID is required", 1)
	}

	if !force {
		fmt.Printf("Are you sure you want to delete server %s? (y/N): ", serverID)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Server deletion cancelled")
			return nil
		}
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	err = client.DeleteServer(serverID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to delete server: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to delete server: %v", err), 1)
	}

	fmt.Printf("Server %s deleted successfully\n", serverID)
	return nil
}