package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

	ao, err := app.Load(appInfo)
	if err != nil {
		return fmt.Errorf("failed to load app output: %w", err)
	}

	printTarget := os.Stdout
	appName := os.Getenv("APP_NAME")

	// print app info
	fmt.Fprintf(printTarget, "Application:\n")
	fmt.Fprintf(printTarget, "  Name: %s\n", appName)

	fmt.Fprintf(printTarget, "\nCD:\n")
	fmt.Fprintf(printTarget, "  URL: %s\n", color.CyanString(ao.CD.DashBoardRef.URL))
	fmt.Fprintf(printTarget, "  Credential:\n")
	fmt.Fprintf(printTarget, "    Username: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Username))
	fmt.Fprintf(printTarget, "    Password: %s\n", color.GreenString(ao.CD.DashBoardRef.Credential.Password))

	// print repos
	fmt.Fprintf(printTarget, "\nRepositories:\n")
	for _, repo := range ao.SCM.Repos {
		fmt.Fprintf(printTarget, "  Name: %s\n", repo.Name)
		fmt.Fprintf(printTarget, "  URL: %s\n\n", color.CyanString(repo.URL))
	}

	fmt.Fprintf(printTarget, "ArgoApps:\n")
	for _, app := range ao.CD.ApplicationRef {
		fmt.Fprintf(printTarget, "  Name: %s\n", app.Name)
	}

	return nil
}
