package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fusioncatalyst/paw/api"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func GenerateAppCodeAction(ctx context.Context, cmd *cli.Command) error {
	// Check if settings file exists
	if _, err := os.Stat("fcsettings.yaml"); os.IsNotExist(err) {
		return cli.Exit("Settings file 'fcsettings.yaml' not found in current directory", 1)
	}

	// Read settings file
	data, err := os.ReadFile("fcsettings.yaml")
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	// Parse settings
	var settings contracts.SettingYAMLFile
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse settings file: %w", err)
	}

	// Get project ID from settings
	if settings.WorkingWithProject == nil {
		return cli.Exit("No project ID specified in settings file. Please run 'paw init-settings-file --working-with-project <project-id>' first", 1)
	}

	// Get app ID from command flags
	appID := cmd.String("app-id")
	if appID == "" {
		return cli.Exit("App ID is required. Please provide it using --app-id flag", 1)
	}

	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}

	// Generate code using API
	code, err := client.GenerateAppCode(appID, settings.CodeGeneration.Language)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Create fusioncat directory if it doesn't exist
	if err := os.MkdirAll("fusioncat", 0755); err != nil {
		return fmt.Errorf("failed to create fusioncat directory: %w", err)
	}

	// Generate filename based on app ID and language
	fileName := fmt.Sprintf("%s.%s", appID, getFileExtension(settings.CodeGeneration.Language))

	// Write the generated code to file
	filePath := filepath.Join("fusioncat", fileName)
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	fmt.Printf("Code generated successfully and saved to %s\n", filePath)
	return nil
}

// getFileExtension returns the appropriate file extension for the given language
func getFileExtension(language string) string {
	switch language {
	case "typescript":
		return "ts"
	case "python":
		return "py"
	case "java":
		return "java"
	case "go":
		return "go"
	default:
		return "txt"
	}
}
