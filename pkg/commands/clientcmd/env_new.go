package clientcmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/commands/clientcmd/state"
)

var (
	envNewCmd = &cobra.Command{
		Use:   "new [NAME]",
		Short: "Create a new environment",
		Long:  "",
		RunE:  envNew,
	}

	envStack string
)

func init() {
	envNewCmd.Flags().StringVar(&envStack, "stack", "", "Specify the name of the stack")
	_ = envNewCmd.MarkFlagRequired("stack")
}

func envNew(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please specify environment name")
	}
	name := args[0]
	_, err := state.InitEnv(name, envStack)
	if err != nil {
		return err
	}

	return nil
}
