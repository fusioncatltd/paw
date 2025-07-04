package tests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/utils"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestProjectsAction(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "paw-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"

	t.Run("List projects without access token", func(t *testing.T) {
		// Create projects command with list subcommand
		projectsCmd := &cli.Command{
			Name: "projects",
			Commands: []*cli.Command{
				{
					Name:   "list",
					Action: actions.ListProjectsAction,
				},
			},
		}

		_, err := utils.CaptureOutputInTests(actions.ListProjectsAction, context.Background(), projectsCmd.Commands[0])
		assert.Contains(t, err.Error(), "401", "Expected 401 error when no access token is set")
	})

	t.Run("Sign up", func(t *testing.T) {
		output, _ := utils.CaptureOutputInTests(actions.SignUpAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "email",
					Value: newUniqueEmail,
				},
				&cli.StringFlag{
					Name:  "password",
					Value: testPassword,
				},
				&cli.BoolFlag{
					Name:  "save-token",
					Value: false,
				},
			},
		})

		token := strings.TrimSpace(string(output))
		if token != "" {
			if err := os.Setenv("FC_ACCESS_TOKEN", token); err != nil {
				t.Fatalf("Failed to store token in environment: %v", err)
			}
		}
	})

	t.Run("List projects with access tokens (but there is not projects yet", func(t *testing.T) {
		// Create projects command with list subcommand
		projectsCmd := &cli.Command{
			Name: "projects",
			Commands: []*cli.Command{
				{
					Name:   "list",
					Action: actions.ListProjectsAction,
				},
			},
		}

		output, _ := utils.CaptureOutputInTests(actions.ListProjectsAction, context.Background(), projectsCmd.Commands[0])
		assert.Contains(t, output, "No projects found", "Expected not to see any projects")
	})

	t.Run("Create a new project", func(t *testing.T) {
		// Create projects command with new subcommand
		projectsCmd := &cli.Command{
			Name: "projects",
			Commands: []*cli.Command{
				{
					Name: "new",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "name",
							Value: "TestProject",
						},
						&cli.StringFlag{
							Name:  "belongs-to",
							Value: "user",
						},
						&cli.StringFlag{
							Name:  "description",
							Value: "Test project description",
						},
						&cli.BoolFlag{
							Name:  "private",
							Value: true,
						},
					},
					Action: actions.CreateNewProjectAction,
				},
			},
		}

		_, err := utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), projectsCmd.Commands[0])
		assert.Nil(t, err)
	})

	t.Run("List projects after creation. TestProject should be in the list", func(t *testing.T) {
		// Create projects command with list subcommand
		projectsCmd := &cli.Command{
			Name: "projects",
			Commands: []*cli.Command{
				{
					Name:   "list",
					Action: actions.ListProjectsAction,
				},
			},
		}

		output, _ := utils.CaptureOutputInTests(actions.ListProjectsAction, context.Background(), projectsCmd.Commands[0])
		assert.Contains(t, output, "TestProject", "Expected to see TestProject in the list of projects")
	})
}
