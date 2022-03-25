package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List items",
	}
)

func init() {
	listCmd.AddCommand(listStackCmd)
}
