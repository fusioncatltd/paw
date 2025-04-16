package actions

import (
	"context"
	"fmt"

	"github.com/fusioncatalyst/paw/api"
	"github.com/urfave/cli/v3"
)

func SignUpAction(_ context.Context, cmd *cli.Command) error {
	email := cmd.String("email")
	password := cmd.String("password")

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

	fmt.Println("Signup successful!")
	return nil
}
