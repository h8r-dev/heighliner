package app

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Output defines the output structure of `hln up` command.
type Output struct {
	ApplicationRef Application `json:"application"`
	Services       []Service   `json:"services,omitempty"`
	CD             CD          `json:"cd,omitempty" yaml:"cd"`
	SCM            SCM         `json:"scm,omitempty" yaml:"scm"`
}

// Application is info about the application itself.
type Application struct {
	Name string `json:"name"`
}

// Service of your app.
type Service struct {
	Name string `json:"name"`
	URL  string `json:"url" yaml:"url"`
	Type string `json:"type" yaml:"type"`
}

// CD now only support argoCD.
type CD struct {
	Provider       string     `json:"provider"`
	Namespace      string     `json:"namespace"`
	Type           string     `json:"type"`
	ApplicationRef []*ArgoApp `json:"applicationRef" yaml:"applicationRef"`
	DashBoardRef   DashBoard  `json:"dashboardRef" yaml:"dashboardRef"`
}

// ArgoApp is argoCD application.
type ArgoApp struct {
	Name        string `json:"name" yaml:"name"`
	Username    string `json:"username,omitempty" yaml:"username"`
	Password    string `json:"password,omitempty" yaml:"password"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
	Infra       string `json:"infra,omitempty" yaml:"infra,omitempty"`
	Prompt      string `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	Type        string `json:"type" yaml:"type"`
	Annotations string `json:"annotations" yaml:"annotations"`
}

// DashBoard information of some component.
type DashBoard struct {
	URL        string     `json:"url"`
	Credential Credential `json:"credential"`
}

// Credential for login to dashboard.
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SCM is source code manager like github.
type SCM struct {
	Provider     string  `json:"provider"`
	Manager      string  `json:"manager"`
	TfProvider   string  `json:"tfProvider" yaml:"tfProvider"`
	Organization string  `json:"organization"`
	Repos        []*Repo `json:"repos"`
}

// Repo is a source code repository.
type Repo struct {
	Name          string        `json:"name"`
	Visibility    string        `json:"visibility"`
	URL           string        `json:"url"`
	TerraformVars TerraformVars `json:"terraformVars" yaml:"terraformVars"`
}

// TerraformVars for deleting repo.
type TerraformVars struct {
	Suffix    string `json:"suffix"`
	Namespace string `json:"namespace"`
}

// Load read and marshal the output yaml file.
// Deprecated
func Load(path string) (*Output, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	output := &Output{}
	err = yaml.Unmarshal(b, output)
	return output, err
}

// ConvertOutputToStatus Convert Output To Status
func (ao *Output) ConvertOutputToStatus() Status {
	s := Status{}
	s.CD.Provider = ao.CD.Provider
	s.CD.URL = ao.CD.DashBoardRef.URL
	s.CD.Username = ao.CD.DashBoardRef.Credential.Username
	s.CD.Password = ao.CD.DashBoardRef.Credential.Password

	s.SCM = ao.SCM

	if len(ao.CD.ApplicationRef) > 0 {
		s.Services = make([]ServiceInfo, 0)
	}

	for _, app := range ao.CD.ApplicationRef {
		a := ServiceInfo{
			Name:     app.Name,
			Type:     app.Type,
			Username: app.Username,
			Password: app.Password,
			URL:      app.URL,
			Prompt:   app.Prompt,
			Infra:    app.Infra,
		}

		s.Services = append(s.Services, a)
	}

	if len(ao.Services) > 0 {
		s.UserServices = make([]UserService, 0)
	}

	for _, service := range ao.Services {
		var u UserService
		u.Service = service

		var repo *Repo
		for _, r := range ao.SCM.Repos {
			if r.Name == service.Name {
				repo = r
				break
			}
		}
		if repo != nil {
			u.Repo = repo
		}
		s.UserServices = append(s.UserServices, u)
	}
	return s
}
