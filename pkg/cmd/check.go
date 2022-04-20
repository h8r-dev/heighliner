package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/nhctl"
)

func newCheckCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if the infrastructures are available",
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		dc, err := dagger.NewDefaultClient(streams)
		if err != nil {
			return err
		}
		if err := dc.Check(); err != nil {
			return err
		}
		if err := nhctl.Check(); err != nil {
			return err
		}
		return util.Exec(streams, nhctl.GetPath(), "version")
	}
	return cmd
}
