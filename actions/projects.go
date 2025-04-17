package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fusioncatalyst/paw/api"
	"github.com/urfave/cli/v3"
)

func ListProjectsAction(_ context.Context, cmd *cli.Command) error {
	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	// Get projects from API
	projects, err := client.ListProjects()
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to list projects: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to list projects: %v", err), 1)
	}

	// If no projects found, show a friendly message
	if len(projects) == 0 {
		fmt.Println("No projects found")
		return nil
	}

	// Format output as JSON for consistency and parsing
	jsonData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format projects: %v", err), 1)
	}

	fmt.Println(string(jsonData))
	return nil
}
