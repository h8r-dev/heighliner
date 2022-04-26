package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/kube"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

// upOptions controls the behavior of up command.
type downOptions struct {
	Dir string

	genericclioptions.IOStreams
}

func (o *downOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
}

func (o *downOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *downOptions) Run() error {
	dClient, err := k8sutil.NewFactory("").DynamicClient()
	if err != nil {
		return err
	}
	return kube.DeleteArgoCDApps(context.Background(), dClient, "argocd", o.IOStreams)
}

func newDownCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &downOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Take down your application",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(cmd, args); err != nil {
				return err
			}
			return o.Run()
		},
	}
	o.BindFlags(cmd.Flags())

	return cmd
}
