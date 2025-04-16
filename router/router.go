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
				Usage:       "paw signup",
				Description: "Create a new user account",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "email",
						Usage:    "The email address of the new user",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Usage:    "The password for the new user",
						Required: true,
					},
				},
				Action: actions.SignUpAction,
			},
		},
	}

	return cmd
}
