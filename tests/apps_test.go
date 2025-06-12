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
	"github.com/fusioncatalyst/paw/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestAppRelatedActions(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := "TestProjectForImport"
	var projectID string // To store the ID of the created project

	validImportFilePathOriginal := "./testfiles/imports/validImport1.yaml"

	t.Run("Set up test", func(t *testing.T) {
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
		assert.Nil(t, err)

		token := strings.TrimSpace(string(output))
		if token != "" {
			if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
				t.Fatalf("Failed to store token in environment: %v", err)
			}
		} else {
			t.Fatal("Signup did not return a token")
		}

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
		assert.Nil(t, err)

		var createdProject api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &createdProject)
		projectID = createdProject.ID // Store the ID for the import step

		assert.NotEmpty(t, projectID, "Project ID should be set before import test")

		output, importErr := utils.CaptureOutputInTests(actions.ImportProjectAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "file",
					Value: validImportFilePathOriginal,
				},
			},
		})
		assert.Empty(t, output, "Output should be empty for valid import")
		assert.Nil(t, importErr, "Import should not return an error")
	})

	t.Run("List apps in project", func(t *testing.T) {
		// Create command with project ID
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		}

		// Execute list apps action
		output, err := utils.CaptureOutputInTests(actions.ListAppsAction, context.Background(), cmd)
		assert.Nil(t, err)

		// Parse the JSON output
		var apps []api.AppAPIResponse
		err = json.Unmarshal([]byte(output), &apps)
		assert.Nil(t, err)

		// Verify the response
		assert.NotNil(t, apps, "Apps list should not be nil")
		assert.GreaterOrEqual(t, len(apps), 1, "Should return more than 1 app")

		// If there are apps, verify their structure
		if len(apps) > 0 {
			for _, app := range apps {
				assert.NotEmpty(t, app.ID, "App ID should not be empty")
				assert.NotEmpty(t, app.Name, "App name should not be empty")
				assert.Equal(t, projectID, app.ProjectID, "App should belong to the test project")
				assert.NotEmpty(t, app.Status, "App status should not be empty")
			}
		}
	})

	t.Run("Create new app", func(t *testing.T) {
		// Create command with required flags
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "TestApp",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Test app description",
				},
			},
		}

		// Execute create app action
		output, err := utils.CaptureOutputInTests(actions.CreateNewAppAction, context.Background(), cmd)
		assert.Nil(t, err)

		// Parse the JSON output
		var createdApp api.AppAPIResponse
		err = json.Unmarshal([]byte(output), &createdApp)
		assert.Nil(t, err)

		// Verify the response
		assert.NotEmpty(t, createdApp.ID, "App ID should not be empty")
		assert.Equal(t, "TestApp", createdApp.Name, "App name should match the input")
		assert.Equal(t, "Test app description", createdApp.Description, "App description should match the input")
		assert.Equal(t, projectID, createdApp.ProjectID, "App should belong to the test project")
		assert.NotEmpty(t, createdApp.Status, "App status should not be empty")

		// Verify the app appears in the list
		listCmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		}

		output, err = utils.CaptureOutputInTests(actions.ListAppsAction, context.Background(), listCmd)
		assert.Nil(t, err)

		var apps []api.AppAPIResponse
		err = json.Unmarshal([]byte(output), &apps)
		assert.Nil(t, err)

		// Find the newly created app in the list
		found := false
		for _, app := range apps {
			if app.ID == createdApp.ID {
				found = true
				assert.Equal(t, "TestApp", app.Name, "App name should match in the list")
				assert.Equal(t, "Test app description", app.Description, "App description should match in the list")
				break
			}
		}
		assert.True(t, found, "Newly created app should be found in the list")
	})
}
