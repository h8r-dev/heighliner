package infra

// Output is the output of initializing infra
type Output struct {
	Dashboard Dashboard `json:"dashboard"`
}

// Dashboard is the dashboard of heighliner
type Dashboard struct {
	Ingress     string      `json:"ingress"`
	Credentials Credentials `json:"credentials"`
}

// Credentials is the credentials of dashboard
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
