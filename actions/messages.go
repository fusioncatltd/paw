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

func ListMessagesAction(ctx context.Context, cmd *cli.Command) error {
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

	// Get list of messages
	messages, err := client.ListMessages(projectID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to list messages: %s", err))
	}

	// Print formatted JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(messages); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}

func CreateMessageAction(ctx context.Context, cmd *cli.Command) error {
	// Get required parameters from command flags
	projectID := cmd.String("project-id")
	name := cmd.String("name")
	description := cmd.String("description")
	schemaID := cmd.String("schema-id")
	schemaVersion := cmd.Int("schema-version")

	if projectID == "" {
		return cli.Exit("Project ID is required. Please provide it using --project-id flag", 1)
	}
	if name == "" {
		return cli.Exit("Message name is required. Please provide it using --name flag", 1)
	}
	if schemaID == "" {
		return cli.Exit("Schema ID is required. Please provide it using --schema-id flag", 1)
	}
	if schemaVersion == 0 {
		return cli.Exit("Schema version is required. Please provide it using --schema-version flag", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	// Create new message
	message, err := client.CreateMessage(projectID, name, description, schemaID, schemaVersion)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create message: %s", err))
	}

	// Print formatted JSON response
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(message); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}
