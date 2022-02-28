package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	stackCmd = &cobra.Command{
		Use:   "stack",
		Short: "Manage stacks",
	}
)

func init() {
	stackCmd.AddCommand(
		stackListCmd,
	)
}
