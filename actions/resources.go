package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/fusioncatalyst/paw/api"
	"github.com/fusioncatalyst/paw/contracts"
)

func ListResourcesAction(ctx context.Context, cmd *cli.Command) error {
	serverID := cmd.String("server-id")

	if serverID == "" {
		return errors.New("Server ID is required")
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize API client: %v", err))
	}

	resources, err := client.ListServerResources(serverID)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(fmt.Sprintf("Failed to list resources: %s", apiErr))
		}
		return errors.New(fmt.Sprintf("Failed to list resources: %v", err))
	}

	jsonData, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to format resources: %v", err))
	}

	fmt.Println(string(jsonData))
	return nil
}

func CreateResourceAction(ctx context.Context, cmd *cli.Command) error {
	serverID := cmd.String("server-id")
	name := cmd.String("name")
	description := cmd.String("description")
	resourceType := cmd.String("type")
	mode := cmd.String("mode")

	if serverID == "" {
		return errors.New("Server ID is required")
	}

	if name == "" {
		return errors.New("Resource name is required")
	}

	if resourceType == "" {
		return errors.New("Resource type is required (topic, exchange, queue, table, endpoint)")
	}

	if mode == "" {
		return errors.New("Resource mode is required (read, write, bind, readwrite)")
	}

	resourceTypeEnum := contracts.ResourceType(resourceType)
	switch resourceTypeEnum {
	case contracts.ResourceTypeKafkaTopic,
		contracts.ResourceTypeExchange,
		contracts.ResourceTypeQueue,
		contracts.ResourceTypeTable,
		contracts.ResourceTypeEndpoint:
	default:
		return errors.New("Invalid resource type. Must be one of: topic, exchange, queue, table, endpoint")
	}

	resourceModeEnum := contracts.ResourceMode(mode)
	switch resourceModeEnum {
	case contracts.ResourceModeRead,
		contracts.ResourceModeWrite,
		contracts.ResourceModeBind,
		contracts.ResourceModeReadWrite:
	default:
		return errors.New("Invalid resource mode. Must be one of: read, write, bind, readwrite")
	}

	client, err := api.NewFCApiClient()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize API client: %v", err))
	}

	resource := contracts.CreateResourceRequest{
		Name:         name,
		Description:  description,
		ResourceType: resourceTypeEnum,
		Mode:         resourceModeEnum,
	}

	newResource, err := client.CreateResource(serverID, resource)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(fmt.Sprintf("Failed to create resource: %s", apiErr))
		}
		return errors.New(fmt.Sprintf("Failed to create resource: %v", err))
	}

	jsonData, err := json.MarshalIndent(newResource, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to format resource: %v", err))
	}

	fmt.Println(string(jsonData))
	return nil
}