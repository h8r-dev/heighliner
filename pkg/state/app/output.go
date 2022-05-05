package app

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

// Output defines the output structure of `hln up` command.
type Output struct {
	ApplicationRef Application `json:"application"`
	Services       []Service   `json:"services,omitempty"`
	CD             CD          `json:"cd,omitempty"`
	SCM            SCM         `json:"scm,omitempty"`
}

// Application is info about the application itself.
type Application struct {
	Name string `json:"name"`
}

// Service of your app.
type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// CD now only support argoCD.
type CD struct {
	Provider       string     `json:"provider"`
	Namespace      string     `json:"namespace"`
	Type           string     `json:"type"`
	ApplicationRef []*ArgoApp `json:"applicationRef"`
	DashBoardRef   DashBoard  `json:"dashboardRef"`
}

// ArgoApp is argoCD application.
type ArgoApp struct {
	Name     string `json:"name"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
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
func (ao *Output) PrettyPrint(streams genericclioptions.IOStreams) error {
	printTarget := streams.Out

	fmt.Fprintln(printTarget, fmt.Sprintf("Heighliner application %s is Ready.", ao.ApplicationRef.Name))
	//fmt.Fprintf(printTarget, "Application:\n")
	//fmt.Fprintf(printTarget, "  Name: %s\n", ao.ApplicationRef.Name)

	fmt.Fprintln(printTarget, fmt.Sprintf("You can access Argocd on %s [Username: %s, Password: %s]", color.CyanString(ao.CD.DashBoardRef.URL),
		color.GreenString(ao.CD.DashBoardRef.Credential.Username), color.GreenString(ao.CD.DashBoardRef.Credential.Password)))

	//fmt.Fprintf(printTarget, "\nCD:\n")
	//fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(ao.CD.DashBoardRef.URL))
	//fmt.Fprintf(printTarget, "  Credential:\n")
	//fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Username))
	//fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Password))

	fmt.Fprintf(printTarget, "\nServices:\n")
	for _, service := range ao.Services {
		fmt.Fprintf(printTarget, "  %s:\n", color.HiBlueString(service.Name))
		fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(service.URL))
	}

	fmt.Fprintf(printTarget, "\nYour repositories in GitHub is:\n")
	for _, repo := range ao.SCM.Repos {
		fmt.Fprintf(printTarget, "  %s:\n", color.HiBlueString(repo.Name))
		fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(repo.URL))
	}

	fmt.Fprintf(printTarget, "\nArgoApps:\n")
	for _, app := range ao.CD.ApplicationRef {
		fmt.Fprintf(printTarget, "  Name: %s\n", app.Name)
		if app.Username != "" {
			fmt.Fprintf(printTarget, "  Credential:\n")
			fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(app.Username))
			fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(app.Password))
		}
	}

	return nil
}
