package cmd

import (
	"github.com/spf13/cobra"
)

// NewListCmd creates and returns the list command of hln
func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List apps or stacks",
	}

	listCmd.AddCommand(newListStacksCmd())

	return listCmd
}
