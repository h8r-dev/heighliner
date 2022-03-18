package commands

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
)

var (
	upCmd = &cobra.Command{
		Use:                "up",
		Short:              "Run an application",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE:               upStack,
	}
)

func upStack(c *cobra.Command, args []string) error {
	newArgs := make([]string, 0)
	newArgs = append(newArgs, "up", "--project", "")
	newArgs = append(newArgs, args...)
	err := util.Exec(
		"dagger",
		newArgs...,
	)
	return err
}
