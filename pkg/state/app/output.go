package app

// Output defines the format of the output from `up` command.
type Output struct {
	Application *Application      `json:"application"`
	Repository  *Repository       `json:"repository,omitempty"`
	Infra       []*InfraComponent `json:"infra,omitempty"`
}

// Application defines the application specific information.
type Application struct {
	Domain  string `json:"domain,omitempty"`
	Ingress string `json:"ingress,omitempty"`
}

// Repository defines the repository specific information.
type Repository struct {
	Backend  string `json:"backend,omitempty"`
	Frontend string `json:"frontend,omitempty"`
	Deploy   string `json:"deploy,omitempty"`
}

// InfraComponent defines the information of an infra component.
type InfraComponent struct {
	Type     string `json:"type"`
	URL      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
