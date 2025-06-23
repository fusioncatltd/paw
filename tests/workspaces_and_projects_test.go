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

func TestWorkspaceAndProjectInteractions(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail-complex-%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	personalProjectName := "mypersonalproject"
	workspaceName := "mytestworkspacreforprojects"
	workspaceProjectName := "myworkspaceproject"
	var workspaceID string

	t.Run("Sign up for complex test", func(t *testing.T) {
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

	t.Run("Create a personal project", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: personalProjectName,
				},
				&cli.StringFlag{
					Name:  "belongs-to",
					Value: "user",
				},
			},
		})
		assert.NoError(t, err)

		var createdProject api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &createdProject)
		assert.NoError(t, err)
		assert.Equal(t, personalProjectName, createdProject.Name)
	})

	t.Run("Create a new workspace", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateWorkspaceAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: workspaceName,
				},
			},
		})
		assert.NoError(t, err)

		var createdWorkspace api.WorkspaceAPIResponse
		err = json.Unmarshal([]byte(output), &createdWorkspace)
		assert.NoError(t, err)
		assert.Equal(t, workspaceName, createdWorkspace.Name)
		workspaceID = createdWorkspace.ID
		assert.NotEmpty(t, workspaceID)
	})

	t.Run("Create a project inside the workspace", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: workspaceProjectName,
				},
				&cli.StringFlag{
					Name:  "belongs-to",
					Value: "workspace",
				},
				&cli.StringFlag{
					Name:  "workspace-id",
					Value: workspaceID,
				},
			},
		})
		assert.NoError(t, err)

		var createdProject api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &createdProject)
		assert.NoError(t, err)
		assert.Equal(t, workspaceProjectName, createdProject.Name)
	})

	t.Run("List all projects and verify count is two", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListProjectsAction, context.Background(), &cli.Command{})
		assert.NoError(t, err)

		var projects []api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &projects)
		assert.NoError(t, err)
		assert.Len(t, projects, 2, "Should have two projects available now")

		projectNames := make([]string, len(projects))
		for i, p := range projects {
			projectNames[i] = p.Name
		}
		assert.Contains(t, projectNames, personalProjectName)
		assert.Contains(t, projectNames, workspaceProjectName)
	})
}
