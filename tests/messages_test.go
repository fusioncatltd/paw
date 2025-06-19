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

func TestMessageActions(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := "TestProjectForMessages"
	var projectID string // To store the ID of the created project
	var schemaID string  // To store the ID of the created schema

	// Get the path to the valid schema file
	schemaFilePath := "testfiles/schemas/validSchema1.json"
	if _, err := os.Stat(schemaFilePath); err != nil {
		t.Fatalf("Failed to find schema file: %v", err)
	}

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

	t.Run("Create a new project", func(t *testing.T) {
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
		assert.Nil(t, err)
		projectID = createdProject.ID
		assert.NotEmpty(t, projectID, "Project ID should be set")
	})

	t.Run("Create a schema for messages", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateSchemaAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "UserSchema",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Schema for user data",
				},
				&cli.StringFlag{
					Name:  "type",
					Value: "jsonschema",
				},
				&cli.StringFlag{
					Name:  "schema-file",
					Value: schemaFilePath,
				},
			},
		})
		assert.Nil(t, err)

		var createdSchema api.SchemaAPIResponse
		err = json.Unmarshal([]byte(output), &createdSchema)
		assert.Nil(t, err)
		schemaID = createdSchema.ID
		assert.NotEmpty(t, schemaID, "Schema ID should be set")
		assert.Equal(t, "UserSchema", createdSchema.Name)
	})

	t.Run("List messages in empty project", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListMessagesAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.Nil(t, err)

		var messages []api.MessageAPIResponse
		err = json.Unmarshal([]byte(output), &messages)
		assert.Nil(t, err)
		assert.Empty(t, messages, "Should return empty list for new project")
	})

	t.Run("Create first message", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateMessageAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "UserCreated",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Message sent when a user is created",
				},
				&cli.StringFlag{
					Name:  "schema-id",
					Value: schemaID,
				},
				&cli.IntFlag{
					Name:  "schema-version",
					Value: 1,
				},
			},
		})
		assert.Nil(t, err)

		var createdMessage api.MessageAPIResponse
		err = json.Unmarshal([]byte(output), &createdMessage)
		assert.Nil(t, err)
		assert.Equal(t, "UserCreated", createdMessage.Name)
		assert.Equal(t, "Message sent when a user is created", createdMessage.Description)
		assert.Equal(t, schemaID, createdMessage.SchemaID)
		assert.Equal(t, 1, createdMessage.SchemaVersion)
	})

	t.Run("Create second message", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateMessageAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "UserUpdated",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Message sent when a user is updated",
				},
				&cli.StringFlag{
					Name:  "schema-id",
					Value: schemaID,
				},
				&cli.IntFlag{
					Name:  "schema-version",
					Value: 1,
				},
			},
		})
		assert.Nil(t, err)

		var createdMessage api.MessageAPIResponse
		err = json.Unmarshal([]byte(output), &createdMessage)
		assert.Nil(t, err)
		assert.Equal(t, "UserUpdated", createdMessage.Name)
		assert.Equal(t, "Message sent when a user is updated", createdMessage.Description)
		assert.Equal(t, schemaID, createdMessage.SchemaID)
		assert.Equal(t, 1, createdMessage.SchemaVersion)
	})

	t.Run("List messages after creation", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListMessagesAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.Nil(t, err)

		var messages []api.MessageAPIResponse
		err = json.Unmarshal([]byte(output), &messages)
		assert.Nil(t, err)
		assert.Len(t, messages, 2, "Should return two messages")

		// Verify message names
		messageNames := make([]string, len(messages))
		for i, msg := range messages {
			messageNames[i] = msg.Name
		}
		assert.Contains(t, messageNames, "UserCreated")
		assert.Contains(t, messageNames, "UserUpdated")
	})

	t.Run("Create message with missing required fields", func(t *testing.T) {
		testCases := []struct {
			name          string
			flags         []cli.Flag
			expectedError string
		}{
			{
				name: "missing project ID",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "Test Message",
					},
					&cli.StringFlag{
						Name:  "schema-id",
						Value: schemaID,
					},
					&cli.IntFlag{
						Name:  "schema-version",
						Value: 1,
					},
				},
				expectedError: "Project ID is required",
			},
			{
				name: "missing message name",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "schema-id",
						Value: schemaID,
					},
					&cli.IntFlag{
						Name:  "schema-version",
						Value: 1,
					},
				},
				expectedError: "Message name is required",
			},
			{
				name: "missing schema ID",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "Test Message",
					},
					&cli.IntFlag{
						Name:  "schema-version",
						Value: 1,
					},
				},
				expectedError: "Schema ID is required",
			},
			{
				name: "missing schema version",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "Test Message",
					},
					&cli.StringFlag{
						Name:  "schema-id",
						Value: schemaID,
					},
				},
				expectedError: "Schema version is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := utils.CaptureOutputInTests(actions.CreateMessageAction, context.Background(), &cli.Command{
					Flags: tc.flags,
				})
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})

	t.Run("Create message with non-existing schema", func(t *testing.T) {
		_, err := utils.CaptureOutputInTests(actions.CreateMessageAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "InvalidMessage",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Message with non-existing schema",
				},
				&cli.StringFlag{
					Name:  "schema-id",
					Value: "non-existing-schema-id",
				},
				&cli.IntFlag{
					Name:  "schema-version",
					Value: 1,
				},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create message")
	})

	t.Run("List messages with invalid project ID", func(t *testing.T) {
		_, err := utils.CaptureOutputInTests(actions.ListMessagesAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: "invalid-project-id",
				},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list messages")
	})
}
