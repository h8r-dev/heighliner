package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hln",
		Short: "Heighliner: Cloud native stacks to develop and deploy your applications",
		Long:  "Heighliner: Cloud native stacks to develop and deploy your applications",
	}
)

func init() {
	rootCmd.AddCommand(stackCmd)
	rootCmd.AddCommand(stackhubCmd)
	rootCmd.AddCommand(upCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
