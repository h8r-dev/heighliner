package app

// Status app status
type Status struct {
	AppName         string // Heighliner app name
	CD              CDInfo
	Apps            []ApplicationInfo
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

// ApplicationInfo app info
type ApplicationInfo struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URL      string `json:"url"`
	Prompt   string `json:"prompt"`
	*Repo
	*Service
}
