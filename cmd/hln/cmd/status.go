package cmd

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

// statusOption controls the behavior of status command.
type statusOption struct {
	KubeconfigPath string
	Namespace      string

	//Kubecli *kubernetes.Clientset

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
	//o.Kubecli = kubecli
	fmt.Println(kubecli)

	// todo: by hxx load from configmap
	// todo: by hxx specify a appName here
	cms, err := kubecli.CoreV1().ConfigMaps(heighlinerNs).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{"type": "app"}).AsSelector().String(),
	})
	if err != nil {
		return err
	}
	if len(cms.Items) == 0 {
		return fmt.Errorf("config map len is 0")
	}

	if len(cms.Items[0].Data) == 0 || cms.Items[0].Data["output.yaml"] == "" {
		return fmt.Errorf("no data in configmap")
	}

	ao := app.Output{}
	err = yaml.Unmarshal([]byte(cms.Items[0].Data["output.yaml"]), &ao)
	if err != nil {
		return err
	}

	//ao, err := app.Load(appInfo)
	//if err != nil {
	//	return fmt.Errorf("failed to load app output: %w", err)
	//}

	if err := ao.PrettyPrint(o.IOStreams); err != nil {
		return err
	}

	return nil
}
