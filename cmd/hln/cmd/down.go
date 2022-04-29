package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/h8r-dev/heighliner/pkg/terraform"
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
	if os.Getenv("GITHUB_TOKEN") == "" {
		return errors.New("please set GITHUB_TOKEN environment variable")
	}
	return nil
}

func (o *downOptions) Run() error {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	pat := os.Getenv("GITHUB_TOKEN")
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

	return deleteRepos(kubeconfig, pat, output.SCM, o.IOStreams)
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
			lg.Info(fmt.Sprintf("argo app %s already deleted", app.Name), zap.NamedError("warn", err))
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

func deleteRepos(kubeconfig, token string, scm app.SCM, streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	os.Setenv("TF_VAR_github_token", token)
	os.Setenv("TF_VAR_organization", scm.Organization)
	tfClient, err := terraform.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	for _, repo := range scm.Repos {
		lg.Info(fmt.Sprintf("delete %s...", repo.Name))
		if err := tfClient.Destroy(terraform.NewApplyOptions(
			terraformDir,
			repo.TerraformVars.Suffix,
			repo.TerraformVars.Namespace,
			kubeconfig,
		)); err != nil {
			return err
		}
	}
	return nil
}
