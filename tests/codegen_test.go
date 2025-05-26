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
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestAppCodegenActions(t *testing.T) {
	// Load .env file
	//if err := godotenv.Load(".env"); err != nil {
	//	t.Fatalf("Error loading .env file: %v", err)
	//}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := "TestProjectForImport"
	var projectID string // To store the ID of the created project

	validImportFilePathOriginal := "./testfiles/imports/validImport1.yaml"

	tempDir, _ := os.MkdirTemp("", "paw-test-*")
	os.Chdir(tempDir)

	t.Run("Set up test", func(t *testing.T) {
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "http://127.0.0.1:8080/",
				},
				&cli.StringFlag{
					Name:  "language",
					Value: "go",
				},
			},
		}
		cmd.Set("server", "http://127.0.0.1:8080/")
		cmd.Set("language", "go")

		output, err := utils.CaptureOutputInTests(actions.InitDefaultSettingsFileAction, context.Background(), cmd)
		assert.NoError(t, err)
		assert.Contains(t, output, "Configuration file 'fcsettings.yaml' has been created")

		output, err = utils.CaptureOutputInTests(actions.SignUpAction, context.Background(), &cli.Command{
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

		output, _ = utils.CaptureOutputInTests(actions.ImportProjectAction, context.Background(), &cli.Command{
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
	})

	//t.Run("Generate code for app", func(t *testing.T) {
	//	// First, get list of apps using CLI command
	//	output, err := utils.CaptureOutputInTests(actions.ListAppsAction, context.Background(), &cli.Command{
	//		Flags: []cli.Flag{
	//			&cli.StringFlag{
	//				Name:  "project-id",
	//				Value: projectID,
	//			},
	//		},
	//	})
	//	assert.Nil(t, err)
	//
	//	var apps []api.AppAPIResponse
	//	err = json.Unmarshal([]byte(output), &apps)
	//	assert.Nil(t, err)
	//	assert.NotEmpty(t, apps, "Should have at least one app after import")
	//
	//	// Pick the first app
	//	appID := apps[0].ID
	//	fmt.Println("Using app ID:", appID)
	//})
}
