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
	"github.com/fusioncatalyst/paw/utils"
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

	// Define test cases for signup, signin, and me action
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
				&cli.BoolFlag{
					Name:  "save-token",
					Value: true, // Save token for subsequent me action test
				},
			},
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "get personal info with valid token",
			action: actions.MeAction,
			flags:  []cli.Flag{}, // MeAction doesn't require flags
			setup: func() error {
				return nil
			},
			expectedErrorToContainString: nil,
		},
		{
			name:   "get personal info without token",
			action: actions.MeAction,
			flags:  []cli.Flag{},
			setup: func() error {
				// Clear the token before this test
				return os.Unsetenv("FC_ACCESS_TOKEN")
			},
			expectedErrorToContainString: &expextintgAuthErrorCode,
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

			// Capture output for me action to verify response format
			var output string
			var err error

			if tt.name == "get personal info with valid token" {
				// Capture and verify the output for me action
				outputBytes, captureErr := utils.CaptureOutputInTests(tt.action, context.Background(), &cli.Command{
					Flags: tt.flags,
				})
				output = string(outputBytes)
				err = captureErr

				// Verify the output contains expected user info fields
				if err == nil {
					if !strings.Contains(output, `"id"`) ||
						!strings.Contains(output, `"handle"`) ||
						!strings.Contains(output, `"status"`) {
						t.Error("Expected output to contain user info fields (id, handle, status)")
					}
				}
			} else {
				// Regular execution for other actions
				err = tt.action(context.Background(), &cli.Command{
					Flags: tt.flags,
				})
			}

			// Check error
			if tt.expectedErrorToContainString != nil {
				if err == nil {
					t.Errorf("Expected error code %v, got nil", tt.expectedErrorToContainString)
				} else if !strings.Contains(err.Error(), *tt.expectedErrorToContainString) {
					t.Errorf("Expected error to contain %v, got %v", *tt.expectedErrorToContainString, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
