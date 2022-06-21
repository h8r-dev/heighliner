package infra

type Output struct {
	Dashboard Dashboard `json:"dashboard"`
}
type Dashboard struct {
	Ingress     string      `json:"ingress"`
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
