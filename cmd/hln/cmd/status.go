package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

// statusOption controls the behavior of status command.
type statusOption struct {
	KubeconfigPath string
	Namespace      string

	Kubecli *kubernetes.Clientset

	genericclioptions.IOStreams
}

func newStatusCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &statusOption{
		IOStreams: streams,
	}
	c := &cobra.Command{
		Use:   "status",
		Short: "Show status of your application",
		Args:  cobra.NoArgs,
		RunE:  o.getStatus,
	}

	if home := homedir.HomeDir(); home != "" {
		c.Flags().StringVar(&o.KubeconfigPath, "", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		c.Flags().StringVar(&o.KubeconfigPath, "", "", "(optional) absolute path to the kubeconfig file")
	}
	c.Flags().StringVar(&o.Namespace, "namespace", "default", "Specify the namespace")

	return c
}

func (o *statusOption) getStatus(c *cobra.Command, args []string) error {
	kubecli, err := k8sutil.NewFactory(o.KubeconfigPath).KubernetesClientSet()
	if err != nil {
		return fmt.Errorf("failed to make kube client: %w", err)
	}
	o.Kubecli = kubecli

	ao, err := app.Load(appInfo)
	if err != nil {
		return fmt.Errorf("failed to load app output: %w", err)
	}

	if err := ao.PrettyPrint(o.IOStreams); err != nil {
		return err
	}

	return nil
}
