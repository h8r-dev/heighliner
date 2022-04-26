package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
}

func newStatusCmd() *cobra.Command {
	o := &statusOption{}
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

	ao, err := loadAppOutput()
	if err != nil {
		return fmt.Errorf("failed to load app output: %w", err)
	}

	printTarget := os.Stdout
	appName := os.Getenv("APP_NAME")

	// print app info
	fmt.Fprintf(printTarget, "Application:\n")
	fmt.Fprintf(printTarget, "  Name: %s\n", appName)
	fmt.Fprintf(printTarget, "  Namespace: %s\n", o.Namespace)
	fmt.Fprintf(printTarget, "  Domain: %s\n", ao.Application.Domain)
	fmt.Fprintf(printTarget, "  IP: %s\n", ao.Application.Ingress)

	// print repos
	fmt.Fprintf(printTarget, "Repositories:\n")
	fmt.Fprintf(printTarget, "  Backend: %s\n", ao.Repository.Backend)
	fmt.Fprintf(printTarget, "  Frontend: %s\n", ao.Repository.Frontend)
	fmt.Fprintf(printTarget, "  Deploy: %s\n", ao.Repository.Deploy)

	// print workload replicas and health
	workload, err := o.loadWorkload(appName, o.Namespace)
	if err != nil {
		return fmt.Errorf("failed to load workload: %w", err)
	}
	fmt.Fprintf(printTarget, "Deployment:\n")
	fmt.Fprintf(printTarget, "  Replicas: %d\n", workload.Spec.Replicas)

	for _, comp := range ao.Infra {
		switch comp.Type {
		case "argoCD":
			fmt.Fprintf(printTarget, "ArgoCD:\n")
			fmt.Fprintf(printTarget, "  Username: %s\n", comp.Username)
			fmt.Fprintf(printTarget, "  Password: %s\n", comp.Password)
		case "grafana":
			fmt.Fprintf(printTarget, "Grafana:\n")
			fmt.Fprintf(printTarget, "  Username: %s\n", comp.Username)
			fmt.Fprintf(printTarget, "  Password: %s\n", comp.Password)
		case "nocalhost":
			fmt.Fprintf(printTarget, "Nocalhost:\n")
			fmt.Fprintf(printTarget, "  Username: %s\n", comp.Username)
			fmt.Fprintf(printTarget, "  Password: %s\n", comp.Password)
		}
	}

	return nil
}

func (o *statusOption) loadWorkload(name, ns string) (*appsv1.Deployment, error) {
	deployment, err := o.Kubecli.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

func loadAppOutput() (*app.Output, error) {
	ao := &app.Output{}
	b, err := ioutil.ReadFile(appInfo)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, ao)
	return ao, err
}
