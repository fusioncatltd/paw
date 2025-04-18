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
			{
				Name:        "projects",
				Usage:       "List all projects",
				Description: "Get information about all projects you have access to",
				Action:      actions.ListProjectsAction,
			},
			{
				Name:        "new-project",
				Usage:       "Create a new project",
				Description: "Create a new project with the specified name and description",
				Action:      actions.CreateNewProjectAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Project name",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "private",
						Usage: "Make the project private",
					},
					&cli.StringFlag{
						Name:     "belongs-to",
						Usage:    "The project can belong to a user or an workspace (which is a group of users)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "workspace-id",
						Usage:    "If project belongs to a workspace, specify the workspace ID",
						Required: false,
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "Optional description of the project",
					},
				},
			},
			{
				Name:        "import-project",
				Usage:       "Import project from dec",
				Description: "Create a new project with the specified name and description",
				Action:      actions.CreateNewProjectAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Usage:    "The id of the project",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "file",
						Usage:    "Path to the file with project definition",
						Required: true,
					},
				},
			},
			{
				Name:        "me",
				Usage:       "paw me",
				Description: "Returns the information about authentication token owner",
				Action:      actions.MeAction,
			},
		},
	}

	return cmd
}
