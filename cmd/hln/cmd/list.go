package cmd

import (
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List apps or stacks",
	}

	// listCmd.AddCommand(newListStacksCmd())
	listCmd.AddCommand(newListAppsCmd())

	return listCmd
}
