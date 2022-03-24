package cmd

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/util"
)

var (
	upCmd = &cobra.Command{
		Use:                "up",
		Short:              "Run an application",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE:               upProj,
	}
)

func upProj(c *cobra.Command, args []string) error {
	// TODO switch to passing args
	newArgs := make([]string, 0)
	newArgs = append(newArgs, "do", "up", "-p", "plans")
	newArgs = append(newArgs, args...)
	err := util.Exec("dagger", newArgs...)
	return err
}
