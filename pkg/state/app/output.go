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
}

// ArgoApp is argoCD application CRD.
type ArgoApp struct {
	Name string `yaml:"name"`
}

// SCM is source code manager like github.
type SCM struct {
	Repos []*Repo `json:"repos"`
}

// Repo is a source code repository.
type Repo struct {
	SecrectSuffix  string `json:"secretsuffix"`
	NameSpace      string `json:"namespace"`
	RepoName       string `json:"repo_name"`
	RepoVisibility string `json:"repo_visibility"`
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
