package app

// UserService user service
type UserService struct {
	Service
	*Repo
}

// Status app status
type Status struct {
	AppName         string // Heighliner app name
	CD              CDInfo
	Services        []ServiceInfo // addon service
	UserServices    []UserService
	SCM             SCM
	TFConfigMapName string
}

// CDInfo CD info
type CDInfo struct {
	Provider string `json:"provider" yaml:"provider"`
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// ServiceInfo service info
type ServiceInfo struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URL      string `json:"url"`
	Infra    string `json:"infra"`
	Prompt   string `json:"prompt"`
	//*Repo
	//*Service
}
