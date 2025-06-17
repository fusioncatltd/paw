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

func ListSchemasAction(ctx context.Context, cmd *cli.Command) error {
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

	// Get list of schemas
	schemas, err := client.ListSchemas(projectID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to list schemas: %s", err))
	}

	// Print formatted JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(schemas); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}

func CreateSchemaAction(ctx context.Context, cmd *cli.Command) error {
	// Get required parameters from command flags
	projectID := cmd.String("project-id")
	name := cmd.String("name")
	description := cmd.String("description")
	schemaType := cmd.String("type")
	schemaFile := cmd.String("schema-file")

	if projectID == "" {
		return cli.Exit("Project ID is required. Please provide it using --project-id flag", 1)
	}
	if name == "" {
		return cli.Exit("Schema name is required. Please provide it using --name flag", 1)
	}
	if schemaType == "" {
		return cli.Exit("Schema type is required. Please provide it using --type flag", 1)
	}
	if schemaFile == "" {
		return cli.Exit("Schema file is required. Please provide it using --schema-file flag", 1)
	}

	// Read schema content from file
	content, err := os.ReadFile(schemaFile)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read schema file: %s", err), 1)
	}
	finalSchemaContent := string(content)

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	// Create new schema
	schema, err := client.CreateSchema(projectID, name, description, schemaType, finalSchemaContent)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create schema: %s", err))
	}

	// Print formatted JSON response
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(schema); err != nil {
		return errors.New(fmt.Sprintf("failed to encode response: %s", err))
	}

	return nil
}
