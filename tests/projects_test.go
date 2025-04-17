package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func TestProjectsAction(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

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

	// Generate unique email for signup
	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"

	// Test cases
	testCases := []struct {
		name                         string
		action                       func(context.Context, *cli.Command) error
		setup                        func() error
		expectedErrorToContainString *string
	}{
		{
			name:   "signup first to get token",
			action: actions.SignUpAction,
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "list projects with valid token",
			action: actions.ListProjectsAction,
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "list projects without token",
			action: actions.ListProjectsAction,
			setup: func() error {
				// Clear the token
				os.Unsetenv("FC_ACCESS_TOKEN")
				return nil
			},
			expectedErrorToContainString: stringPtr("401"), // Expecting unauthorized error
		},
	}

	// Capture original stdout
	oldStdout := os.Stdout

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			var cmd *cli.Command
			if tt.name == "signup first to get token" {
				// Create a pipe to capture stdout
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("Failed to create pipe: %v", err)
				}
				os.Stdout = w

				// Setup signup command
				cmd = &cli.Command{
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "email",
							Value: newUniqueEmail,
						},
						&cli.StringFlag{
							Name:  "password",
							Value: testPassword,
						},
						&cli.BoolFlag{
							Name:  "save-token",
							Value: true,
						},
					},
				}

				// Run the action
				err = tt.action(context.Background(), cmd)

				// Close the write end of the pipe
				w.Close()

				// Read the captured output
				output, err := ioutil.ReadAll(r)
				if err != nil {
					t.Fatalf("Failed to read captured output: %v", err)
				}

				// Restore stdout
				os.Stdout = oldStdout

				// Store the captured token in environment
				token := strings.TrimSpace(string(output))
				if token != "" {
					if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
						t.Fatalf("Failed to store token in environment: %v", err)
					}
				}
			} else {
				// Setup projects command
				cmd = &cli.Command{}
				// Run the action
				err = tt.action(context.Background(), cmd)
			}

			// Check error
			if tt.expectedErrorToContainString != nil {
				if err == nil {
					t.Errorf("Expected error containing %v, got nil", *tt.expectedErrorToContainString)
				} else if !strings.Contains(err.Error(), *tt.expectedErrorToContainString) {
					t.Errorf("Expected error to contain %v, got %v", *tt.expectedErrorToContainString, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
