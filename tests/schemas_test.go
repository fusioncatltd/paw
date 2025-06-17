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

func TestSchemaActions(t *testing.T) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("testmail%s@testmail.com", currentTimestamp)
	testPassword := "password123"
	projectName := "TestProjectForSchemas"
	var projectID string // To store the ID of the created project

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

	t.Run("List schemas in empty project", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListSchemasAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.Nil(t, err)

		var schemas []api.SchemaAPIResponse
		err = json.Unmarshal([]byte(output), &schemas)
		assert.Nil(t, err)
		assert.Empty(t, schemas, "Should return empty list for new project")
	})

	t.Run("Create schema with file", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.CreateSchemaAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
				&cli.StringFlag{
					Name:  "name",
					Value: "PersonSchema",
				},
				&cli.StringFlag{
					Name:  "description",
					Value: "Schema for a person with first name, last name, and age",
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
		assert.Equal(t, "PersonSchema", createdSchema.Name)
		assert.Equal(t, "Schema for a person with first name, last name, and age", createdSchema.Description)
	})

	t.Run("List schemas after creation", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListSchemasAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.Nil(t, err)

		var schemas []api.SchemaAPIResponse
		err = json.Unmarshal([]byte(output), &schemas)
		assert.Nil(t, err)
		assert.Len(t, schemas, 1, "Should return one schema")

		// Verify schema name
		assert.Equal(t, "PersonSchema", schemas[0].Name)
	})

	// Store schema ID for update tests
	var schemaID string
	t.Run("Get schema ID for update tests", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListSchemasAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "project-id",
					Value: projectID,
				},
			},
		})
		assert.Nil(t, err)

		var schemas []api.SchemaAPIResponse
		err = json.Unmarshal([]byte(output), &schemas)
		assert.Nil(t, err)
		assert.Len(t, schemas, 1, "Should return one schema")
		schemaID = schemas[0].ID
		assert.NotEmpty(t, schemaID, "Schema ID should be set")
	})

	// Create a temporary file with updated schema content
	updatedSchemaContent := `{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"title": "Person",
		"type": "object",
		"properties": {
			"firstName": {
				"type": "string",
				"description": "The person's first name."
			},
			"lastName": {
				"type": "string",
				"description": "The person's last name."
			},
			"age": {
				"description": "Age in years which must be a non-negative integer.",
				"type": "integer",
				"minimum": 0
			},
			"email": {
				"type": "string",
				"format": "email",
				"description": "The person's email address."
			}
		},
		"required": ["firstName", "lastName", "email"]
	}`

	tempFile, err := os.CreateTemp("", "updated-schema-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if err := os.WriteFile(tempFile.Name(), []byte(updatedSchemaContent), 0644); err != nil {
		t.Fatalf("Failed to write updated schema content: %v", err)
	}

	t.Run("Update schema content", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.UpdateSchemaAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "schema-id",
					Value: schemaID,
				},
				&cli.StringFlag{
					Name:  "schema-file",
					Value: tempFile.Name(),
				},
			},
		})
		assert.Nil(t, err)

		var updatedSchema api.SchemaAPIResponse
		err = json.Unmarshal([]byte(output), &updatedSchema)
		assert.Nil(t, err)
		assert.Equal(t, schemaID, updatedSchema.ID)
		assert.Equal(t, "PersonSchema", updatedSchema.Name)
		assert.Equal(t, "Schema for a person with first name, last name, and age", updatedSchema.Description)
	})

	t.Run("Update schema with missing required fields", func(t *testing.T) {
		testCases := []struct {
			name          string
			flags         []cli.Flag
			expectedError string
		}{
			{
				name: "missing schema ID",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "schema-file",
						Value: tempFile.Name(),
					},
				},
				expectedError: "Schema ID is required",
			},
			{
				name: "missing schema file",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "schema-id",
						Value: schemaID,
					},
				},
				expectedError: "Schema file is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := utils.CaptureOutputInTests(actions.UpdateSchemaAction, context.Background(), &cli.Command{
					Flags: tc.flags,
				})
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})

	t.Run("Create schema with missing required fields", func(t *testing.T) {
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
						Value: "Test Schema",
					},
					&cli.StringFlag{
						Name:  "type",
						Value: "json",
					},
					&cli.StringFlag{
						Name:  "schema-file",
						Value: schemaFilePath,
					},
				},
				expectedError: "Project ID is required",
			},
			{
				name: "missing schema name",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "type",
						Value: "json",
					},
					&cli.StringFlag{
						Name:  "schema-file",
						Value: schemaFilePath,
					},
				},
				expectedError: "Schema name is required",
			},
			{
				name: "missing schema type",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "Test Schema",
					},
					&cli.StringFlag{
						Name:  "schema-file",
						Value: schemaFilePath,
					},
				},
				expectedError: "Schema type is required",
			},
			{
				name: "missing schema file",
				flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project-id",
						Value: projectID,
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "Test Schema",
					},
					&cli.StringFlag{
						Name:  "type",
						Value: "json",
					},
				},
				expectedError: "Schema file is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := utils.CaptureOutputInTests(actions.CreateSchemaAction, context.Background(), &cli.Command{
					Flags: tc.flags,
				})
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})
}
