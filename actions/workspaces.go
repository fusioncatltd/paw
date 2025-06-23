package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fusioncatalyst/paw/api"
	"github.com/urfave/cli/v3"
)

func ListWorkspacesAction(ctx context.Context, cmd *cli.Command) error {
	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	// Get list of workspaces
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to list workspaces: %s", err))
	}

	// Print formatted JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(workspaces); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}

func CreateWorkspaceAction(ctx context.Context, cmd *cli.Command) error {
	// Get required parameters from command flags
	name := cmd.String("name")
	description := cmd.String("description")

	if name == "" {
		return cli.Exit("Workspace name is required. Please provide it using --name flag", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	// Create new workspace
	workspace, err := client.CreateWorkspace(name, description)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create workspace: %s", err))
	}

	// Print formatted JSON response
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(workspace); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}
