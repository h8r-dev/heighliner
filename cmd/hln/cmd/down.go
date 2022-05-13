package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	"github.com/h8r-dev/heighliner/internal/k8sfactory"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/h8r-dev/heighliner/pkg/terraform"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

// upOptions controls the behavior of up command.
type downOptions struct {
	Dir              string
	IsDeletePackages bool

	genericclioptions.IOStreams
}

func (o *downOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.BoolVar(&o.IsDeletePackages, "delete-packages", false, "Delete packages")
}

func (o *downOptions) Validate(cmd *cobra.Command, args []string) error {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return errors.New("please set GITHUB_TOKEN environment variable")
	}
	return nil
}

func (o *downOptions) Run(appName string) error {
	kubeconfig := k8sutil.GetKubeConfigPath()
	pat := os.Getenv("GITHUB_TOKEN")

	state, err := getStateInSpecificBackend()
	if err != nil {
		return err
	}
	output, err := state.LoadOutput(appName)
	if err != nil {
		return err
	}

	dClient, err := k8sfactory.GetDefaultFactory().DynamicClient()
	if err != nil {
		return err
	}
	if err := deleteArgoCDApps(context.Background(), dClient, output.CD, o.IOStreams); err != nil {
		return err
	}

	if o.IsDeletePackages {
		if err := deletePackages(pat, output.SCM, o.IOStreams); err != nil {
			return err
		}
	}

	return deleteRepos(appName, kubeconfig, pat, output.SCM, o.IOStreams)
}

func newDownCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &downOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "down [appName]",
		Short: "Take down your application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(cmd, args); err != nil {
				return err
			}
			return o.Run(args[0])
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
	lg.Info(fmt.Sprintf("delete app %s...", name))
	return argoApp.Delete(ctx, name, metav1.DeleteOptions{})
}

func deleteRepos(appName, kubeconfig, token string, scm app.SCM, streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	if err := os.Setenv("TF_VAR_github_token", token); err != nil {
		return err
	}
	if err := os.Setenv("TF_VAR_organization", scm.Organization); err != nil {
		return err
	}
	tfClient, err := terraform.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	for _, repo := range scm.Repos {
		lg.Info(fmt.Sprintf("delete repo %s...", repo.Name))
		repoDir := filepath.Join(terraformDir, repo.Name)
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			return err
		}
		tfContent, err := GetTFProvider(appName)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(repoDir, "provider.tf"), []byte(tfContent), 0644)
		if err != nil {
			return err
		}
		if err := tfClient.Destroy(terraform.NewApplyOptions(
			repoDir,
			repo.TerraformVars.Suffix,
			repo.TerraformVars.Namespace,
			kubeconfig,
		)); err != nil {
			return err
		}
	}
	return nil
}

func deletePackages(token string, scm app.SCM, streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)

	// set GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	gClient := github.NewClient(tc)

	pkgType := "container"
	for _, repo := range scm.Repos {

		if _, _, err := gClient.Users.GetPackage(ctx, scm.Organization, pkgType, repo.Name); err != nil {
			if strings.Contains(err.Error(), "404 Package not found.") {
				continue
			}
			return err
		}

		lg.Info(fmt.Sprintf("delete package %s...", repo.Name))
		if _, err := gClient.Users.DeletePackage(ctx, scm.Organization, pkgType, repo.Name); err != nil {
			return err
		}
	}
	return nil
}
