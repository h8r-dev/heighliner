package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"
)

// Output defines the format of the output from `up` command.
type Output struct {
	ApplicationRef Application `json:"application"`
	Services       []Service   `json:"services,omitempty"`
	CD             CD          `json:"cd,omitempty"`
	SCM            SCM         `json:"scm,omitempty"`
}

type Application struct {
	Name string `json:"name"`
}

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

func (ao *Output) PrettyPrint(streams genericclioptions.IOStreams) error {
	printTarget := streams.Out

	fmt.Fprintf(printTarget, "Application:\n")
	fmt.Fprintf(printTarget, "  Name: %s\n", ao.ApplicationRef.Name)

	fmt.Fprintf(printTarget, "\nServices:\n")
	for _, service := range ao.Services {
		fmt.Fprintf(printTarget, "  %s:\n", color.HiBlueString(service.Name))
		fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(service.URL))
	}

	fmt.Fprintf(printTarget, "\nCD:\n")
	fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(ao.CD.DashBoardRef.URL))
	fmt.Fprintf(printTarget, "  Credential:\n")
	fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Username))
	fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Password))

	fmt.Fprintf(printTarget, "\nRepositories:\n")
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
