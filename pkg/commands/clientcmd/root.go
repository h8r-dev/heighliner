package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hln",
		Short: "Heighliner: Cloud native best practices to build and deploy your applications",
		Long:  "Heighliner: Cloud native best practices to build and deploy your applications",
	}
)

func init() {
	rootCmd.AddCommand(
		initCmd,
		stackCmd,
		stackhubCmd,
		envCmd,
		configCmd,
		secretCmd,
		upCmd,
	)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
