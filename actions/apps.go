package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fusioncatalyst/paw/api"
	"github.com/urfave/cli/v3"
)

func ListAppsAction(ctx context.Context, cmd *cli.Command) error {
	// Get project ID from command flags
	projectID := cmd.String("project-id")
	if projectID == "" {
		return cli.Exit("Project ID is required. Please provide it using --project-id flag", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}

	// Get list of apps
	apps, err := client.ListApps(projectID)
	if err != nil {
		return fmt.Errorf("failed to list apps: %w", err)
	}

	// Print formatted JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(apps); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
