package app

// Status app status
type Status struct {
	AppName         string // Heighliner app name
	Cd              CdInfo
	Apps            []ApplicationInfo
	SCM             SCM
	TfConfigMapName string
}

// CdInfo Cd info
type CdInfo struct {
	Provider string `json:"provider" yaml:"provider"`
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// ApplicationInfo app info
type ApplicationInfo struct {
	Name     string
	Type     string
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	*Repo
	*Service
}