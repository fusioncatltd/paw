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

func TestWorkspaceActions(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	workspaceName1 := "testworkspace1"
	workspaceName2 := "testworkspace2"
	var workspaceID1 string

	t.Run("Sign up", func(t *testing.T) {
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
		assert.NoError(t, err)

		token := strings.TrimSpace(string(output))
		if token != "" {
			if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
				t.Fatalf("Failed to store token in environment: %v", err)
			}
		} else {
			t.Fatal("Signup did not return a token")
		}
	})

	t.Run("List workspaces for new user", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListWorkspacesAction, context.Background(), &cli.Command{})
		assert.NoError(t, err)

		var workspaces []api.UserWorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &workspaces)
		assert.NoError(t, err)
		assert.Empty(t, workspaces, "Should return an empty list for a new user")
	})

	t.Run("Create first workspace", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateWorkspaceAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: workspaceName1,
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "This is the first test workspace.",
				},
			},
		})
		assert.NoError(t, err)

		var createdWorkspace api.WorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &createdWorkspace)
		assert.NoError(t, err)
		assert.Equal(t, workspaceName1, createdWorkspace.Name)
		assert.Equal(t, "This is the first test workspace.", createdWorkspace.Description)
		assert.NotEmpty(t, createdWorkspace.ID, "Workspace ID should be set")
		workspaceID1 = createdWorkspace.ID
	})

	t.Run("List workspaces after creating one", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListWorkspacesAction, context.Background(), &cli.Command{})
		assert.NoError(t, err)

		var workspaces []api.UserWorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &workspaces)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 1, "Should return one workspace")
		assert.Equal(t, workspaceID1, workspaces[0].Workspace.ID)
		assert.Equal(t, workspaceName1, workspaces[0].Workspace.Name)
	})

	t.Run("Create second workspace", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateWorkspaceAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: workspaceName2,
				},
			},
		})
		assert.NoError(t, err)

		var createdWorkspace api.WorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &createdWorkspace)
		assert.NoError(t, err)
		assert.Equal(t, workspaceName2, createdWorkspace.Name)
		assert.Empty(t, createdWorkspace.Description, "Description should be empty as it was not provided")
		assert.NotEmpty(t, createdWorkspace.ID, "Workspace ID should be set")
	})

	t.Run("List workspaces after creating two", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListWorkspacesAction, context.Background(), &cli.Command{})
		assert.NoError(t, err)

		var workspaces []api.UserWorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &workspaces)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 2, "Should return two workspaces")

		// Verify workspace names
		workspaceNames := make([]string, len(workspaces))
		for i, ws := range workspaces {
			workspaceNames[i] = ws.Workspace.Name
		}
		assert.Contains(t, workspaceNames, workspaceName1)
		assert.Contains(t, workspaceNames, workspaceName2)
	})

	t.Run("Create workspace with missing name", func(t *testing.T) {
		_, err := utils.CaptureOutputInTests(actions.CreateWorkspaceAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Workspace name is required")
	})
}
