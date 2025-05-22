package tests

import (
	"context"
	"github.com/fusioncatalyst/paw/actions"
	"os"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestInitDefaultSettingsFileAction(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "paw-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	tests := []struct {
		name          string
		setup         func() error
		expectedError error
		expectedFile  string
	}{
		{
			name: "successful file creation",
			setup: func() error {
				return nil
			},
			expectedError: nil,
			expectedFile:  "fcsettings.yaml",
		},
		{
			name: "file already exists",
			setup: func() error {
				// Create the file first
				file, err := os.Create("fcsettings.yaml")
				if err != nil {
					return err
				}
				return file.Close()
			},
			expectedError: cli.Exit("File 'fcsettings.yaml' already exists in current directory", 1),
			expectedFile:  "fcsettings.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Run the action
			err := actions.InitDefaultSettingsFileAction(context.Background(), &cli.Command{})

			// Check error
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check file existence if expected
			if tt.expectedFile != "" && tt.expectedError == nil {
				if _, err := os.Stat(tt.expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s to exist, but it doesn't", tt.expectedFile)
				}
			}
		})
	}
}
