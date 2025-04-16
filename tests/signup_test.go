package tests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func TestSignUpAction(t *testing.T) {
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

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)

	tests := []struct {
		name          string
		flags         []cli.Flag
		setup         func() error
		expectedError error
	}{
		{
			name: "successful signup",
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: "password123",
				},
			},
			setup: func() error {
				return nil
			},
			expectedError: nil,
		},
		//{
		//	name:      "missing email",
		//	arguments: []string{"", "password123"},
		//	setup: func() error {
		//		return nil
		//	},
		//	expectedError: fmt.Errorf("both email and password are required"),
		//},
		//{
		//	name:      "missing password",
		//	arguments: []string{"test@example.com", ""},
		//	setup: func() error {
		//		return nil
		//	},
		//	expectedError: fmt.Errorf("both email and password are required"),
		//},
		//{
		//	name:      "no arguments",
		//	arguments: []string{},
		//	setup: func() error {
		//		return nil
		//	},
		//	expectedError: fmt.Errorf("both email and password are required"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Create CLI command with actual arguments
			//cmd := &cli.Command{
			//	Name:      "signup",
			//	Arguments: tt.arguments, // Directly set the arguments
			//}

			// Run the action
			err := actions.SignUpAction(context.Background(), &cli.Command{
				Flags: tt.flags,
			})

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
		})
	}
}
