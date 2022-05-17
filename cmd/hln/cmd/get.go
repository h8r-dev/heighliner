package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newGetCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newGetIngressCmd(streams))

	return cmd
}
