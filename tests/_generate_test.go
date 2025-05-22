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

	"github.com/fusioncatalyst/paw/utils"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/fusioncatalyst/paw/api"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestGenerateCodeAction(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := "TestProjectForGenerate"
	var projectID string // To store the ID of the created project

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
				&cli.BoolFlag{
					Name:  "save-token",
					Value: false,
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
	})

	t.Run("Create a new project to generate code for", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), &cli.Command{
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
		projectID = createdProject.ID // Store the ID for the generate step
	})

	t.Run("Generate code for app", func(t *testing.T) {
		assert.NotEmpty(t, projectID, "Project ID should be set before generate test")

		// Create context with required values
		ctx := context.WithValue(context.Background(), "project-id", projectID)
		ctx = context.WithValue(ctx, "app-id", "test-app")

		// Create generate command
		generateCmd := &cli.Command{
			Name: "generate",
		}

		output, err := utils.CaptureOutputInTests(actions.GenerateCodeAction, ctx, generateCmd)
		assert.Nil(t, err, "Generate code failed")

		// Optional: Verify the output contains expected info
		var generatedProject api.ProjectAPIResponse
		err = json.Unmarshal([]byte(output), &generatedProject)
		assert.Nil(t, err, "Failed to parse generated project response")
		assert.Equal(t, projectID, generatedProject.ID, "Generated project ID should match the original")
		assert.Equal(t, projectName, generatedProject.Name, "Generated project Name should match the original")
	})

	t.Run("Generate code with invalid app ID", func(t *testing.T) {
		assert.NotEmpty(t, projectID, "Project ID should be set before generate test")

		// Create context with required values
		ctx := context.WithValue(context.Background(), "project-id", projectID)
		ctx = context.WithValue(ctx, "app-id", "invalid-app-id")

		// Create generate command
		generateCmd := &cli.Command{
			Name: "generate",
		}

		_, err := utils.CaptureOutputInTests(actions.GenerateCodeAction, ctx, generateCmd)
		assert.NotNil(t, err, "Expected an error for invalid app ID")
		assert.Contains(t, err.Error(), "App not found", "Error message should indicate app not found")
	})
}
