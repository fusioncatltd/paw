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

func ListAppsAction(ctx context.Context, cmd *cli.Command) error {
	// Get project ID from command flags
	projectID := cmd.String("project-id")
	if projectID == "" {
		return cli.Exit("Project ID is required. Please provide it using --project-id flag", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	// Get list of apps
	apps, err := client.ListApps(projectID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to list apps: %s", err))
	}

	// Print formatted JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(apps); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}
