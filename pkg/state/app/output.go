package app

import (
	"os"

	"sigs.k8s.io/yaml"
)

// Output defines the format of the output from `up` command.
type Output struct {
	CD  CD  `json:"cd"`
	SCM SCM `json:"scm"`
}

// CD is information about argoCD.
type CD struct {
	Provider       string     `json:"provider"`
	Namespace      string     `json:"namespace"`
	Type           string     `json:"type"`
	ApplicationRef []*ArgoApp `json:"applicationRef"`
	DashBoardRef   DashBoard  `json:"dashboardRef"`
}

// ArgoApp is argoCD application CRD.
type ArgoApp struct {
	Name     string `json:"name"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// DashBoard of some component.
type DashBoard struct {
	URL        string     `json:"url"`
	Credential Credential `json:"credential"`
}

// Credential ifor login.
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SCM is source code manager like github.
type SCM struct {
	Provider     string  `json:"provider"`
	Manager      string  `json:"manager"`
	TfProvider   string  `json:"tfProvider"`
	Organization string  `json:"organization"`
	Repos        []*Repo `json:"repos"`
}

// Repo is a source code repository.
type Repo struct {
	Name          string        `json:"name"`
	Visibility    string        `json:"visibility"`
	URL           string        `json:"url"`
	TerraformVars TerraformVars `json:"terraformVars"`
}

// TerraformVars for deleting repo.
type TerraformVars struct {
	Suffix    string `json:"suffix"`
	Namespace string `json:"namespace"`
}

// Load read and marshal the output yaml file.
func Load(path string) (*Output, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	output := &Output{}
	err = yaml.Unmarshal(b, output)
	return output, err
}
