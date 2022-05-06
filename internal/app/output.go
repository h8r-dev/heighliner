package app

import (
	"gopkg.in/yaml.v3"
	"os"
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

// PrettyPrint format and print the output.
//func (ao *Output) PrettyPrint(streams genericclioptions.IOStreams) error {
//	printTarget := streams.Out
//
//	fmt.Fprintln(printTarget, fmt.Sprintf("Heighliner application %s is Ready.", ao.ApplicationRef.Name))
//	//fmt.Fprintf(printTarget, "Application:\n")
//	//fmt.Fprintf(printTarget, "  Name: %s\n", ao.ApplicationRef.Name)
//
//	fmt.Fprintln(printTarget, fmt.Sprintf("You can access Argocd on %s [Username: %s, Password: %s]", color.CyanString(ao.CD.DashBoardRef.URL),
//		color.GreenString(ao.CD.DashBoardRef.Credential.Username), color.GreenString(ao.CD.DashBoardRef.Credential.Password)))
//
//	//fmt.Fprintf(printTarget, "\nCD:\n")
//	//fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(ao.CD.DashBoardRef.URL))
//	//fmt.Fprintf(printTarget, "  Credential:\n")
//	//fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Username))
//	//fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Password))
//
//	fmt.Fprintf(printTarget, "\nServices:\n")
//	for _, service := range ao.Services {
//		fmt.Fprintf(printTarget, "  %s:\n", color.HiBlueString(service.Name))
//		fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(service.URL))
//	}
//
//	fmt.Fprintf(printTarget, "\nYour repositories in GitHub is:\n")
//	for _, repo := range ao.SCM.Repos {
//		fmt.Fprintf(printTarget, "  %s:\n", color.HiBlueString(repo.Name))
//		fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(repo.URL))
//	}
//
//	fmt.Fprintf(printTarget, "\nArgoApps:\n")
//	for _, app := range ao.CD.ApplicationRef {
//		fmt.Fprintf(printTarget, "  Name: %s\n", app.Name)
//		if app.Username != "" {
//			fmt.Fprintf(printTarget, "  Credential:\n")
//			fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(app.Username))
//			fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(app.Password))
//		}
//	}
//
//	return nil
//}

func (ao *Output) ConvertOutputToStatus() Status {
	s := Status{}
	s.Cd.Provider = ao.CD.Provider
	s.Cd.URL = ao.CD.DashBoardRef.URL
	s.Cd.Username = ao.CD.DashBoardRef.Credential.Username
	s.Cd.Password = ao.CD.DashBoardRef.Credential.Password

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
				svc = &service
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
