package actions

import (
	"context"
	"fmt"
	"os"

	"github.com/fusioncatalyst/paw/api"
	"github.com/urfave/cli/v3"
)

func SignUpAction(_ context.Context, cmd *cli.Command) error {
	email := cmd.String("email")
	password := cmd.String("password")
	saveToken := cmd.Bool("save-token")

	if email == "" || password == "" {
		return fmt.Errorf("both email and password are required")
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}

	if err := client.SignUp(email, password); err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return cli.Exit(fmt.Sprintf("Signup failed: %s", apiErr), 1)
		}
		return cli.Exit(fmt.Sprintf("Signup failed: %v", err), 1)
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
