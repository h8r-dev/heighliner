package cmd

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/dagger"
	"github.com/h8r-dev/heighliner/pkg/util/nhctl"
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if the infrastructures are available",
	}
	cmd.Run = func(c *cobra.Command, args []string) {
		lg := logger.New()
		err := dagger.Check()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to install dagger")
		}
		err = nhctl.Check()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to install nhctl")
		}
		err = util.Exec(dagger.GetPath(), "version")
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute dagger version")
		}
		err = util.Exec(nhctl.GetPath(), "version")
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute nhctl version")
		}
	}
	return cmd
}
