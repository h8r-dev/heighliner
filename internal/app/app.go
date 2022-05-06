package app

type Status struct {
	Cd   CdInfo
	Apps []ApplicationInfo
	SCM  SCM
}

type CdInfo struct {
	Provider string `json:"provider" yaml:"provider"`
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type ApplicationInfo struct {
	Name     string
	Type     string
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	*Repo
	*Service
}
