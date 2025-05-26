package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

func CreateNewProjectAction(_ context.Context, cmd *cli.Command) error {
	// Get required parameters from flags
	name := cmd.String("name")
	belongsTo := cmd.String("belongs-to")
	isPrivate := cmd.Bool("private")

	// Get optional parameters
	description := cmd.String("description")
	workspaceID := cmd.String("workspace-id")

	// Validate belongs-to parameter
	if belongsTo != "user" && belongsTo != "workspace" {
		return cli.Exit("belongs-to must be either 'user' or 'workspace'", 1)
	}

	// If belongs-to is workspace, workspace-id is required
	if belongsTo == "workspace" && workspaceID == "" {
		return cli.Exit("workspace-id is required when belongs-to is set to 'workspace'", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	// Set the created_by_id based on belongs-to
	createdByID := ""
	if belongsTo == "workspace" {
		createdByID = workspaceID
	} else {
		userData, personalInfoError := client.GetPersonalInfo()
		if personalInfoError != nil {
			return personalInfoError
		}
		createdByID = userData.ID
	}

	// Create project using API client
	project, err := client.CreateProject(
		name,
		description,
		belongsTo,
		createdByID,
		isPrivate,
	)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to create project: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to create project: %v", err), 1)
	}

	// Format output as JSON for consistency
	jsonData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format project data: %v", err), 1)
	}

	fmt.Println(string(jsonData))
	return nil
}

func ImportProjectAction(ctx context.Context, cmd *cli.Command) error {
	// Get project ID and file path from context
	projectID := cmd.String("project-id")
	filePath := cmd.String("file")

	// Verify file exists before making API call
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("File not found: %s", filePath), 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err.Error()), 1)
	}

	// Import project using API client
	err = client.ImportProject(projectID, filePath)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to import project: %s", apiErr.Error()), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to import project: %v", err.Error()), 1)
	}

	return nil
}

func GenerateCodeAction(ctx context.Context, cmd *cli.Command) error {
	// Get project ID and app ID from context
	projectID := ctx.Value("project-id").(string)
	appID := ctx.Value("app-id").(string)

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	// Generate code using API client
	project, err := client.GenerateCode(projectID, appID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to generate code: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to generate code: %v", err), 1)
	}

	// Format output as JSON for consistency
	jsonData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format project data: %v", err), 1)
	}

	fmt.Println(string(jsonData))
	return nil
}
