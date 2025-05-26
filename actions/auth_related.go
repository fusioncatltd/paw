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

func SignInAction(_ context.Context, cmd *cli.Command) error {
	email := cmd.String("email")
	password := cmd.String("password")
	saveToken := cmd.Bool("save-token")

	if email == "" || password == "" {
		return errors.New("both email and password are required")
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize API client: %s", err))
	}

	if err := client.SignIn(email, password); err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Sign in failed: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Sign in failed: %v", err), 1)
	}

	// Get the authorization token from the client
	token := client.GetAuthorization()

	// Always print the token to stdout
	fmt.Println(token)

	// If save-token flag is set, store in environment variable
	if saveToken {
		if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to save token to environment: %v", err), 1)
		}
		fmt.Fprintln(os.Stderr, "Token saved to FC_ACCESS_TOKEN environment variable")
	}

	return nil
}

func SignUpAction(_ context.Context, cmd *cli.Command) error {
	email := cmd.String("email")
	password := cmd.String("password")
	saveToken := cmd.Bool("save-token")

	if email == "" || password == "" {
		return cli.Exit("both email and password are required", 1)
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit("failed to initialize API client: "+err.Error(), 1)
	}

	err = client.SignUp(email, password)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(fmt.Sprintf("Signup failed: %s", apiErr))
		}
		return errors.New(fmt.Sprintf("Signup failed: %v", err))
	}

	// Get the authorization token from the client
	token := client.GetAuthorization()

	// Always print the token to stdout
	fmt.Println(token)

	// If save-token flag is set, store in environment variable
	if saveToken {
		if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to save token to environment: %v", err), 1)
		}
	}

	return nil
}

func MeAction(_ context.Context, cmd *cli.Command) error {
	// Initialize API client
	client, err := api.NewFCApiClient()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to initialize API client: %v", err), 1)
	}

	// Get user info from API
	userInfo, err := client.GetPersonalInfo()
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Failed to get user info: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Failed to get user info: %v", err), 1)
	}

	// Format output as JSON for consistency
	jsonData, err := json.MarshalIndent(userInfo, "", "  ")
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to format user info: %v", err), 1)
	}

	fmt.Println(string(jsonData))
	return nil
}
