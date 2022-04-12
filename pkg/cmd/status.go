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

func newStatusCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "status",
		Short: "Show status of your application",
		Args:  cobra.NoArgs,
		RunE:  getStatus,
	}

	if home := homedir.HomeDir(); home != "" {
		c.Flags().StringP("kubeconfig", "", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		c.Flags().StringP("kubeconfig", "", "", "(optional) absolute path to the kubeconfig file")
	}

	return c
}

func getStatus(c *cobra.Command, args []string) error {
	kubeconfigPath := c.Flags().Lookup("kubeconfig").Value.String()
	kubecli, err := k8sutil.MakeKubeClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to make kube client: %w", err)
	}

	ao, err := loadAppOutput()
	if err != nil {
		return fmt.Errorf("failed to load app output: %w", err)
	}

	printTarget := os.Stdout
	appName := os.Getenv("APP_NAME")

	// print app info
	fmt.Fprintf(printTarget, "Application:\n")
	fmt.Fprintf(printTarget, "  Name: %s\n", appName)
	fmt.Fprintf(printTarget, "  Namespace: %s\n", ao.Application.Namespace)
	fmt.Fprintf(printTarget, "  Domain: %s\n", ao.Application.Domain)
	fmt.Fprintf(printTarget, "  IP: %s\n", ao.Application.Ingress)

	// print workload replicas and health
	workload, err := loadWorkload(kubecli, appName, ao.Application.Namespace)
	if err != nil {
		return fmt.Errorf("failed to load workload: %w", err)
	}
	fmt.Fprintf(printTarget, "Deployment:\n")
	fmt.Fprintf(printTarget, "  Replicas: %d\n", workload.Spec.Replicas)

	// print repos

	// print build pipelines and images

	// print argocd app

	// print db
	return nil
}

func loadWorkload(kubecli *kubernetes.Clientset, name, ns string) (*appsv1.Deployment, error) {
	deployment, err := kubecli.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

func loadAppOutput() (*app.Output, error) {
	ao := &app.Output{}
	b, err := ioutil.ReadFile(appOutputPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, ao)
	return ao, err
}
