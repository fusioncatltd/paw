package contracts

import "time"

type ResourceType string
type ResourceMode string

const (
	ResourceTypeKafkaTopic ResourceType = "topic"
	ResourceTypeExchange   ResourceType = "exchange"
	ResourceTypeQueue      ResourceType = "queue"
	ResourceTypeTable      ResourceType = "table"
	ResourceTypeEndpoint   ResourceType = "endpoint"
)

const (
	ResourceModeRead      ResourceMode = "read"
	ResourceModeWrite     ResourceMode = "write"
	ResourceModeBind      ResourceMode = "bind"
	ResourceModeReadWrite ResourceMode = "readwrite"
)

type CreateResourceRequest struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ResourceType ResourceType `json:"resource_type"`
	Mode         ResourceMode `json:"mode"`
	ServerID     string       `json:"server_id"`
}

type ResourceResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	ResourceType      string    `json:"resource_type"`
	Mode              string    `json:"mode"`
	ServerID          string    `json:"server_id"`
	ProjectID         string    `json:"project_id"`
	CreatedByUserID   string    `json:"created_by_user_id"`
	CreatedByUserName string    `json:"created_by_user_name"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
