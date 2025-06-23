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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "server",
						Usage: "Server URL (e.g., https://api.fusioncat.dev)",
					},
					&cli.StringFlag{
						Name:  "language",
						Usage: "Target language (typescript, python, java, go)",
					},
					&cli.StringFlag{
						Name:  "working-with-project",
						Usage: "Connect settings file to a specific project (must be a valid UUID)",
					},
				},
				Action: actions.InitDefaultSettingsFileAction,
			},
			{
				Name:        "auth",
				Usage:       "Authentication commands",
				Description: "Sign up, sign in, and manage authentication",
				Commands: []*cli.Command{
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
						Name:        "me",
						Usage:       "paw me",
						Description: "Returns the information about authentication token owner",
						Action:      actions.MeAction,
					},
				},
			},
			{
				Name:        "codegen",
				Usage:       "Generate code from project definitions",
				Description: "List, create, and manage projects",
				Commands: []*cli.Command{
					{
						Name:        "app",
						Usage:       "Generate code for a specific app",
						Description: "Generate code for a specific application in the project",
						Action:      actions.GenerateAppCodeAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "app-id",
								Usage:    "The ID of the application to generate code for",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "language",
								Usage:    "The target language for code generation (typescript, python, java, go)",
								Required: false,
							},
						},
					},
				},
			},
			{
				Name:        "apps",
				Usage:       "Manage apps",
				Description: "List, create, and manage apps in projects",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List all apps in projects",
						Description: "Get information about all apps in projects you have access to",
						Action:      actions.ListAppsAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to operate on",
								Required: true,
							},
						},
					},
					{
						Name:        "new",
						Usage:       "Create a new app",
						Description: "Create a new app in the specified project",
						Action:      actions.CreateNewAppAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to create the app in",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Name of the app",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Usage:    "Description of the app",
								Required: false,
							},
						},
					},
				},
			},
			{
				Name:        "schemas",
				Usage:       "Manage schemas",
				Description: "List, create, and manage schemas in projects",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List all schemas in a project",
						Description: "Get information about all schemas in a project you have access to",
						Action:      actions.ListSchemasAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to operate on",
								Required: true,
							},
						},
					},
					{
						Name:        "new",
						Usage:       "Create a new schema",
						Description: "Create a new schema in the specified project",
						Action:      actions.CreateSchemaAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to create the schema in",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Name of the schema",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Usage:    "Description of the schema",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "type",
								Usage:    "Type of the schema (e.g., jsonschema)",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "schema-file",
								Usage:    "Path to a file containing the schema",
								Required: true,
							},
						},
					},
					{
						Name:        "update",
						Usage:       "Update an existing schema",
						Description: "Update an existing schema. Note: Only schema content can be updated. Schema name, type, and description cannot be changed.",
						Action:      actions.UpdateSchemaAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "schema-id",
								Usage:    "The ID of the schema to update",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "schema-file",
								Usage:    "Path to a file containing the updated schema content",
								Required: true,
							},
						},
					},
					{
						Name:        "versions",
						Usage:       "List all versions of a schema",
						Description: "Get information about all versions of a schema, including who made the changes and when",
						Action:      actions.ListSchemaVersionsAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "schema-id",
								Usage:    "The ID of the schema to list versions for",
								Required: true,
							},
						},
					},
					{
						Name:        "get-version",
						Usage:       "Get a specific version of a schema",
						Description: "Retrieve a specific version of a schema by its schema ID and version ID",
						Action:      actions.GetSchemaVersionAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "schema-id",
								Usage:    "The ID of the schema",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "version-id",
								Usage:    "The version ID of the schema",
								Required: true,
							},
						},
					},
				},
			},
			{
				Name:        "messages",
				Usage:       "Manage messages",
				Description: "List, create, and manage messages in projects",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List all messages in a project",
						Description: "Get information about all messages in a project you have access to",
						Action:      actions.ListMessagesAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to operate on",
								Required: true,
							},
						},
					},
					{
						Name:        "new",
						Usage:       "Create a new message",
						Description: "Create a new message in the specified project",
						Action:      actions.CreateMessageAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to create the message in",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Name of the message",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Usage:    "Description of the message",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "schema-id",
								Usage:    "The ID of the schema this message is based on",
								Required: true,
							},
							&cli.IntFlag{
								Name:     "schema-version",
								Usage:    "The version of the schema to use",
								Required: true,
							},
						},
					},
				},
			},
			{
				Name:        "projects",
				Usage:       "Manage projects",
				Description: "List, create, and manage projects",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List all projects",
						Description: "Get information about all projects you have access to",
						Action:      actions.ListProjectsAction,
					},
					{
						Name:        "new",
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
						Name:        "import",
						Usage:       "Import project from file",
						Description: "Import project definition from a file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "file",
								Usage:    "Path to the file with project definition",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to operate on",
								Required: true,
							},
						},
						Action: actions.ImportProjectAction,
					},
					{
						Name:        "generate",
						Usage:       "Generate code for project",
						Description: "Generate code for a specific application in the project",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "app-id",
								Usage:    "The ID of the application to generate code for",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "project-id",
								Usage:    "The ID of the project to operate on",
								Required: true,
							},
						},
						Action: actions.GenerateCodeAction,
					},
				},
			},
			{
				Name:        "workspaces",
				Usage:       "Manage workspaces",
				Description: "List and create workspaces",
				Commands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List all workspaces",
						Description: "Get information about all workspaces you have access to",
						Action:      actions.ListWorkspacesAction,
					},
					{
						Name:        "new",
						Usage:       "Create a new workspace",
						Description: "Create a new workspace with the specified name and description",
						Action:      actions.CreateWorkspaceAction,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Workspace name",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "description",
								Usage: "Optional description of the workspace",
							},
						},
					},
				},
			},
		},
	}

	return cmd
}
