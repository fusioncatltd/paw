package tests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func TestSignUpAndSignInAction(t *testing.T) {
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
	testPassword := "password123"

	expextintgAuthErrorCode := "401"

	// Define test cases for both signup and signin
	testCases := []struct {
		name                         string
		action                       func(context.Context, *cli.Command) error
		flags                        []cli.Flag
		setup                        func() error
		expectedErrorToContainString *string
	}{
		{
			name:   "successful signup",
			action: actions.SignUpAction,
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: testPassword,
				},
			},
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "successful signin with created account",
			action: actions.SignInAction,
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: testPassword,
				},
			},
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "signin with wrong password",
			action: actions.SignInAction,
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: "wrongpassword",
				},
			},
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: &expextintgAuthErrorCode,
		},
		{
			name:   "signin with non-existent email",
			action: actions.SignInAction,
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: "nonexistent@testmail.com",
				},
				&cli.StringFlag{
					Name:  "password",
					Value: testPassword,
				},
			},
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: &expextintgAuthErrorCode,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Run the action
			err := tt.action(context.Background(), &cli.Command{
				Flags: tt.flags,
			})

			// Check error
			if tt.expectedErrorToContainString != nil {
				if err == nil {
					t.Errorf("Expected error code %v, got nil", tt.expectedErrorToContainString)
				} else if !strings.Contains(err.Error(), *tt.expectedErrorToContainString) {
					t.Errorf("Expected error to contain  %v, got %v", *tt.expectedErrorToContainString, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
