package cmd

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/state"
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
	t := state.NewTemp()
	if err := t.Detect(); err != nil {
		return err
	}
	newArgs := make([]string, 0)
	newArgs = append(newArgs, "do", "up", "-p", "./plans")
	newArgs = append(newArgs, args...)
	return util.Exec("dagger", newArgs...)
}
