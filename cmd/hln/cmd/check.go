package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/checker"
)

func newCheckCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if the infrastructures are available",
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return checker.Check(streams)
	}
	// Shadow the root PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}
	return cmd
}
