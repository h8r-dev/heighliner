package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newListCmd(streams genericclioptions.IOStreams) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List apps or stacks",
	}

	listCmd.AddCommand(newListStacksCmd(streams))
	listCmd.AddCommand(newListAppsCmd(streams))

	return listCmd
}
