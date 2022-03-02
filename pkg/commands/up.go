package commands

import (
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		Args:  cobra.NoArgs,
		RunE:  upStack,
	}
)

func upStack(c *cobra.Command, args []string) error {
	err := util.Exec(
		"dagger",
		"--project", "",
		"up",
	)
	return err
}
