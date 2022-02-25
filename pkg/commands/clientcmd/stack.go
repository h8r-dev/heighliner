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

func init() {
	stackCmd.AddCommand(
		stackPullCmd,
		stackShowCmd,
		stackInitCmd,
		stackListCmd,
		stackInputCmd,
	)
}
