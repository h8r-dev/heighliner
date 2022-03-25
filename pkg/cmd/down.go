package cmd

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/util"
)

var (
	downCmd = &cobra.Command{
		Use:                "down",
		Short:              "Remove the project and clear",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE:               downProj,
	}
)

func downProj(c *cobra.Command, args []string) error {
	// TODO switch to passing args
	newArgs := make([]string, 0)
	newArgs = append(newArgs, "do", "down", "-p", "plans")
	newArgs = append(newArgs, args...)
	return util.Exec("dagger", newArgs...)
}
