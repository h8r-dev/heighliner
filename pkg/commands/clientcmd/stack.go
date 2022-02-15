package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	stackCmd = &cobra.Command{
		Use:   "stack",
		Short: "Manage stacks",
		Long:  "",
	}
)

var defaultStacks = []*stack{{
	Name:        "sample",
	Version:     "1.0.0",
	Description: "This is an example stack",
	Inputs: []*inputSchema{{
		Name:        "metadata.name",
		Type:        "string",
		Required:    true,
		Description: "The name of the application",
	}, {
		Name:        "push.target",
		Type:        "config",
		Required:    true,
		Description: "The target to push image to",
	}, {
		Name:        "kubeconfig",
		Type:        "secret",
		Required:    true,
		Description: "The kubeconfig to connect to the k8s cluster",
	}},
}}

type inputSchema struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

type stack struct {
	Name        string
	Version     string
	Description string
	Inputs      []*inputSchema
}

func init() {
	stackCmd.AddCommand(stackListCmd)
	stackCmd.AddCommand(stackPullCmd)
	stackCmd.AddCommand(stackShowCmd)
}
