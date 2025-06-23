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
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Resources   []ServerResource  `json:"resources,omitempty"`
	Binds       []ServerBind      `json:"binds,omitempty"`
	ProjectID   string            `json:"project_id,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
}

type CreateServerRequest struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Resources   []ServerResource  `json:"resources,omitempty"`
	Binds       []ServerBind      `json:"binds,omitempty"`
	ProjectID   string            `json:"project_id"`
}

type UpdateServerRequest struct {
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
	Resources   []ServerResource  `json:"resources,omitempty"`
	Binds       []ServerBind      `json:"binds,omitempty"`
}

type ServersListResponse struct {
	Servers []Server `json:"servers"`
	Total   int      `json:"total"`
}