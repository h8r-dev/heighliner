package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/nhctl"
)

func newCheckCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if the infrastructures are available",
	}
	cmd.Run = func(c *cobra.Command, args []string) {
		lg := logger.New()
		dc, err := dagger.NewDefaultClient(streams)
		if err != nil {
			lg.Fatal().Err(err)
		}
		err = dc.Check()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to check dagger version")
		}
		err = nhctl.Check()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to install nhctl")
		}
		err = util.Exec(streams, nhctl.GetPath(), "version")
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute nhctl version")
		}
	}
	return cmd
}
