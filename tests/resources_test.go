package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fusioncatalyst/paw/contracts"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/fusioncatalyst/paw/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestResourceManagement(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	currentTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newUniqueEmail := fmt.Sprintf("resource_test%s@testmail.com", currentTimestamp)
	testPassword := "password123"

	var projectID string
	var serverID string

	output, err := utils.CaptureOutputInTests(actions.SignUpAction, context.Background(), &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "email", Value: newUniqueEmail},
			&cli.StringFlag{Name: "password", Value: testPassword},
		},
	})
	assert.NoError(t, err)
	token := strings.TrimSpace(string(output))
	if token != "" {
		os.Setenv("FC_ACCESS_TOKEN", token)
	}

	projectName := fmt.Sprintf("resourcetestproject%s", currentTimestamp[:10])
	output, err = utils.CaptureOutputInTests(actions.CreateNewProjectAction, context.Background(), &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: projectName},
			&cli.StringFlag{Name: "belongs-to", Value: "user"},
			&cli.BoolFlag{Name: "private", Value: true},
		},
	})
	assert.Nil(t, err)

	var createdProject struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err = json.Unmarshal([]byte(output), &createdProject)
	assert.Nil(t, err)
	assert.Equal(t, projectName, createdProject.Name)
	projectID = createdProject.ID

	t.Run("Create server for resources", func(t *testing.T) {
		serverName := fmt.Sprintf("testserver%s", currentTimestamp[:10])
		serverType := "async+kafka"
		serverDescription := "Test Kafka server for resource testing"

		output, err := utils.CaptureOutputInTests(actions.CreateServer, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "project-id", Value: projectID},
				&cli.StringFlag{Name: "name", Value: serverName},
				&cli.StringFlag{Name: "type", Value: serverType},
				&cli.StringFlag{Name: "description", Value: serverDescription},
			},
		})
		assert.Nil(t, err, "Failed to create server")

		var server contracts.ServerResponse
		err = json.Unmarshal([]byte(output), &server)
		assert.Nil(t, err, "Failed to parse server response")
		assert.Equal(t, serverName, server.Name)
		assert.Equal(t, serverType, server.Protocol)
		serverID = server.ID

		fmt.Println("Server created with ID:", serverID)
	})

	t.Run("Create resource", func(t *testing.T) {
		resourceName := fmt.Sprintf("testtopic%s", currentTimestamp[:10])
		resourceType := "topic"
		resourceMode := "readwrite"
		resourceDescription := "Test Kafka topic resource"

		output, err := utils.CaptureOutputInTests(actions.CreateResourceAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "server-id", Value: serverID},
				&cli.StringFlag{Name: "name", Value: resourceName},
				&cli.StringFlag{Name: "type", Value: resourceType},
				&cli.StringFlag{Name: "mode", Value: resourceMode},
				&cli.StringFlag{Name: "description", Value: resourceDescription},
			},
		})
		assert.Nil(t, err, "Failed to create resource")

		var resource contracts.ResourceResponse
		err = json.Unmarshal([]byte(output), &resource)
		assert.Nil(t, err, "Failed to parse resource response")
		assert.Equal(t, resourceName, resource.Name)
		assert.Equal(t, resourceType, resource.ResourceType)
		assert.Equal(t, resourceMode, resource.Mode)
		assert.Equal(t, resourceDescription, resource.Description)
		assert.Equal(t, serverID, resource.ServerID)
	})

	t.Run("List resources", func(t *testing.T) {
		output, err := utils.CaptureOutputInTests(actions.ListResourcesAction, context.Background(), &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "server-id", Value: serverID},
			},
		})
		assert.Nil(t, err, "Failed to list resources")

		var resources []contracts.ResourceResponse
		err = json.Unmarshal([]byte(output), &resources)
		assert.Nil(t, err, "Failed to parse resources list")
		assert.GreaterOrEqual(t, len(resources), 1, "Should have at least one resource")
	})
}
