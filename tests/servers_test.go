package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/fusioncatalyst/paw/api"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/fusioncatalyst/paw/router"
	"github.com/fusioncatalyst/paw/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestServerActions(t *testing.T) {
	err := godotenv.Load(".env")
	assert.NoError(t, err, "Error loading .env file")

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testservers%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := fmt.Sprintf("TestServerProject%s", currentTimestamp)

	var accessToken string
	var projectID string
	var serverIDs []string

	cmd := router.GetCLIRouter()
	var serversCmd *cli.Command
	for _, c := range cmd.Commands {
		if c.Name == "servers" {
			serversCmd = c
			break
		}
	}
	assert.NotNil(t, serversCmd, "Expected to find 'servers' command")

	validServerTypes := []string{
		"async+kafka",
		"async+amqp",
		"async+mqtt",
		"async+db",
		"async+webhook",
	}

	t.Run("Set up test environment", func(t *testing.T) {
		// Sign up
		output, err := utils.CaptureOutputInTests(actions.SignUpAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: testPassword,
				},
			},
		})
		assert.NoError(t, err, "Failed to sign up")

		// Get the token from signup output
		accessToken = strings.TrimSpace(output)
		assert.NotEmpty(t, accessToken, "Access token should not be empty")

		os.Setenv("FC_ACCESS_TOKEN", accessToken)

		// Create project attached to user
		output, err = utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: projectName,
				},
				&cli.StringFlag{
					Name:  "belongs-to",
					Value: "user",
				},
				&cli.BoolFlag{
					Name:  "private",
					Value: true,
				},
			},
		})
		assert.NoError(t, err)

		var project api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &project)
		assert.NoError(t, err)
		projectID = project.ID
		assert.NotEmpty(t, projectID, "Project ID should not be empty")
	})

	t.Run("Create servers with all valid types", func(t *testing.T) {
		for _, serverType := range validServerTypes {
			serverName := fmt.Sprintf("Test%sServer%s", strings.ReplaceAll(serverType, "+", ""), currentTimestamp)

			output, err := utils.CaptureOutputInTests(actions.CreateServer, context.Background(), &cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "name",
						Value: serverName,
					},
					&cli.StringFlag{
						Name:  "type",
						Value: serverType,
					},
					&cli.StringFlag{
						Name:  "description",
						Value: fmt.Sprintf("Test %s server for automated tests", serverType),
					},
				},
			})
			assert.NoError(t, err, "Failed to create server of type %s", serverType)

			var server contracts.Server
			err = json.Unmarshal([]byte(output), &server)
			assert.NoError(t, err)
			assert.Equal(t, serverName, server.Name)
			assert.Equal(t, serverType, server.Protocol) // API returns 'protocol' not 'type'
			assert.NotEmpty(t, server.ID)
			serverIDs = append(serverIDs, server.ID)
		}
	})

	t.Run("List servers in project", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListServers, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.NoError(t, err)

		var response contracts.ServersListResponse
		err = json.Unmarshal([]byte(output), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response.Servers), len(serverIDs))

		// All servers in response should belong to our project
		for _, server := range response.Servers {
			assert.Equal(t, projectID, server.ProjectID)
		}
	})
	// Error cases
	t.Run("Create server with invalid type", func(t *testing.T) {
		_, err := utils.CaptureOutputInTests(actions.CreateServer, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "InvalidServer",
				},
				&cli.StringFlag{
					Name:  "type",
					Value: "invalid+type",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Should fail",
				},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "422")
	})
	t.Run("Create server without required fields", func(t *testing.T) {
		testCases := []struct {
			name          string
			flags         []cli.Flag
			expectedError string
		}{
			{
				name: "missing name",
				flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Value: projectID},
					&cli.StringFlag{Name: "type", Value: "async+kafka"},
					&cli.StringFlag{Name: "description", Value: "Test"},
				},
				expectedError: "Server name is required",
			},
			{
				name: "missing type",
				flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Value: projectID},
					&cli.StringFlag{Name: "name", Value: "Test"},
					&cli.StringFlag{Name: "description", Value: "Test"},
				},
				expectedError: "Server type is required",
			},
			{
				name: "missing description",
				flags: []cli.Flag{
					&cli.StringFlag{Name: "project-id", Value: projectID},
					&cli.StringFlag{Name: "name", Value: "Test"},
					&cli.StringFlag{Name: "type", Value: "async+kafka"},
				},
				expectedError: "Server description is required",
			},
			{
				name: "missing project ID",
				flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Value: "Test"},
					&cli.StringFlag{Name: "type", Value: "async+kafka"},
					&cli.StringFlag{Name: "description", Value: "Test"},
				},
				expectedError: "Project ID is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := utils.CaptureOutputInTests(actions.CreateServer, context.Background(), &cli.Command{
					Flags: tc.flags,
				})
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})
	t.Run("List servers without project ID", func(t *testing.T) {
		_, err := utils.CaptureOutputInTests(actions.ListServers, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "project-id", Value: ""},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "project ID is required")
	})

	t.Run("Operations without access token", func(t *testing.T) {
		// Remove access token
		os.Unsetenv("FC_ACCESS_TOKEN")

		_, err := utils.CaptureOutputInTests(actions.ListServers, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "project-id", Value: projectID},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "401")

		// Restore token for cleanup
		os.Setenv("FC_ACCESS_TOKEN", accessToken)
	})
}
