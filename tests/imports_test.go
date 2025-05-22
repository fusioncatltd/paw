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

func TestImportProjectAction(t *testing.T) {
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

	t.Run("Create a new project to import into", func(t *testing.T) {
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
		projectID = createdProject.ID // Store the ID for the import step
	})

	t.Run("Import project with valid file", func(t *testing.T) {
		assert.NotEmpty(t, projectID, "Project ID should be set before import test")

		// Create context with required values
		ctx := context.WithValue(context.Background(), "project-id", projectID)
		ctx = context.WithValue(ctx, "file", validImportFilePathOriginal)

		// Create import command
		importCmd := &cli.Command{
			Name: "import",
		}

		output, _ := utils.CaptureOutputInTests(actions.ImportProjectAction, ctx, importCmd)
		assert.Empty(t, output, "Output should be empty for valid import")

	})

	t.Run("Import project with non-existent file", func(t *testing.T) {
		assert.NotEmpty(t, projectID, "Project ID should be set before import test")

		// Create context with required values
		ctx := context.WithValue(context.Background(), "project-id", projectID)
		ctx = context.WithValue(ctx, "file", "non-existent-file.yaml")

		// Create import command
		importCmd := &cli.Command{
			Name: "import",
		}

		_, err := utils.CaptureOutputInTests(actions.ImportProjectAction, ctx, importCmd)
		assert.NotNil(t, err, "Expected an error for non-existent file")
		assert.Contains(t, err.Error(), "File not found", "Error message should indicate file not found")
	})
}
