package cmd

import (
	"github.com/spf13/cobra"
)

// NewListCmd creates and returns the list command of hln
func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Parent command of list stack",
	}

	listCmd.AddCommand(NewListStackCmd())

	return listCmd
}
