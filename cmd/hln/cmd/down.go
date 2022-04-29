package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/state/app"
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
	output, err := app.Load(appInfo)
	if err != nil {
		return err
	}
	dClient, err := k8sutil.NewFactory("").DynamicClient()
	if err != nil {
		return err
	}
	if err := deleteArgoCDApps(context.Background(), dClient, output.CD, o.IOStreams); err != nil {
		return err
	}
	return nil
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

func deleteArgoCDApps(ctx context.Context,
	dClient dynamic.Interface, cd app.CD,
	streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	for _, app := range cd.ApplicationRef {
		if err := patchFinalizerAndDelete(ctx, dClient, cd.Namespace, app.Name, streams); err != nil {
			lg.Info(color.HiYellowString("skip %s and continue", app.Name), zap.NamedError("warn", err))
		}
	}
	return nil
}

func patchFinalizerAndDelete(ctx context.Context,
	client dynamic.Interface,
	namespace, name string,
	streams genericclioptions.IOStreams) error {
	const argoCDFinalizerRaw = `{"metadata": {"finalizers": ["resources-finalizer.argocd.argoproj.io"]}}`
	lg := logger.New(streams)
	var argoAppResource = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}
	argoApp := client.Resource(argoAppResource).Namespace(namespace)
	_, err := argoApp.Patch(ctx, name, types.MergePatchType, []byte(argoCDFinalizerRaw), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	lg.Info(fmt.Sprintf("patch finalizer to app %s", name))
	return argoApp.Delete(ctx, name, metav1.DeleteOptions{})
}
