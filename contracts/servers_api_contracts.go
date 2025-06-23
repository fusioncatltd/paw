package contracts

type ServerResource struct {
	Name         string `json:"name"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	ResourceName string `json:"resource_name,omitempty"`
}

type ServerBind struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	RoutingKey  string `json:"routing_key,omitempty"`
}

type Server struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	Protocol           string `json:"protocol"`
	Description        string `json:"description"`
	Status             string `json:"status,omitempty"`
	ProjectID          string `json:"project_id,omitempty"`
	UserID             string `json:"user_id,omitempty"`
	CreatedByUserName  string `json:"created_by_user_name,omitempty"`
}

type CreateServerRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
}

type ServersListResponse struct {
	Servers []Server `json:"servers"`
	Total   int      `json:"total"`
}