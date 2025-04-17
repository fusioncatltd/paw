package router

import (
	"github.com/fusioncatalyst/paw/actions"
	"github.com/urfave/cli/v3"
)

func GetCLIRouter() *cli.Command {
	cmd := &cli.Command{
		Name:        "paw",
		Version:     "0.1.0",
		Description: "An official fusioncat CLI",
		Arguments:   cli.AnyArguments,
		Commands: []*cli.Command{
			{
				Name:        "init-settings-file",
				Usage:       "paw init-settings-file",
				Description: "Create a settings file with default values",
				Action:      actions.InitDefaultSettingsFileAction,
			},
			{
				Name:        "signup",
				Usage:       "Create a new user account",
				Description: "Sign up for a new account using email and password",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "email",
						Usage:    "User's email address",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Usage:    "User's password",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "save-token",
						Usage: "Save the authorization token to FC_ACCESS_TOKEN environment variable",
					},
				},
				Action: actions.SignUpAction,
			},
			{
				Name:        "signin",
				Usage:       "Sign in to existing account",
				Description: "Sign in using email and password",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "email",
						Usage:    "User's email address",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Usage:    "User's password",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "save-token",
						Usage: "Save the authorization token to FC_ACCESS_TOKEN environment variable",
					},
				},
				Action: actions.SignInAction,
			},
		},
	}

	return cmd
}
