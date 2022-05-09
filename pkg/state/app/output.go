package app

import (
	"os"

	"gopkg.in/yaml.v3"
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
	Name     string `json:"name"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Type     string
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
	TerraformVars TerraformVars `json:"terraformVars"`
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
		s.Apps = make([]ApplicationInfo, 0)
	}

	for _, app := range ao.CD.ApplicationRef {
		a := ApplicationInfo{
			Name:     app.Name,
			Type:     app.Type,
			Username: app.Username,
			Password: app.Password,
		}

		var repo *Repo
		for _, r := range ao.SCM.Repos {
			if r.Name == app.Name {
				repo = r
				break
			}
		}
		if repo != nil {
			a.Repo = repo
		}

		var svc *Service
		for _, service := range ao.Services {
			if service.Name == app.Name {
				sv := service
				svc = &sv
				break
			}
		}
		if svc != nil {
			a.Service = svc
		}
		s.Apps = append(s.Apps, a)
	}

	return s
}
