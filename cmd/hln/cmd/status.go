package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/h8r-dev/heighliner/internal/app"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newStatusCmd(streams genericclioptions.IOStreams) *cobra.Command {
	c := &cobra.Command{
		Use:   "status",
		Short: "Show status of your application",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getStatus()
		},
	}

	return c
}

// For hln down
func getTfProvider() (string, error) {
	s, err := getAppStatus()
	if err != nil {
		return "", err
	}

	if s.TfConfigMapName == "" {
		return "", fmt.Errorf("No tf provider config map? ")
	}

	cli, err := getDefaultClientSet()
	if err != nil {
		return "", err
	}

	cm, err := cli.CoreV1().ConfigMaps(heighlinerNs).Get(context.TODO(), s.TfConfigMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(cm.Data) == 0 || cm.Data[tfProviderConfigMapKey] == "" {
		return "", fmt.Errorf("No data found in tf provider configmap ")
	}
	return cm.Data[tfProviderConfigMapKey], nil
}

// Get Heighliner application status from k8s configmap
func getAppStatus() (*app.Status, error) {
	kubecli, err := getDefaultClientSet()
	if err != nil {
		return nil, fmt.Errorf("failed to make kube client: %w", err)
	}

	// todo: by hxx specify a appName here
	cms, err := kubecli.CoreV1().ConfigMaps(heighlinerNs).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{configTypeKey: "heighliner"}).AsSelector().String(),
	})
	if err != nil {
		return nil, err
	}
	if len(cms.Items) == 0 {
		return nil, fmt.Errorf("config map len is 0")
	}

	if len(cms.Items[0].Data) == 0 || cms.Items[0].Data["output.yaml"] == "" {
		return nil, fmt.Errorf("no data in configmap")
	}

	appName := cms.Items[0].Name

	ao := app.Output{}
	err = yaml.Unmarshal([]byte(cms.Items[0].Data["output.yaml"]), &ao)
	if err != nil {
		return nil, err
	}

	s := ao.ConvertOutputToStatus()
	if cms.Items[0].Data[tfProviderConfigMapKey] != "" {
		s.TfConfigMapName = cms.Items[0].Data[tfProviderConfigMapKey]
	}
	s.AppName = appName
	return &s, nil
}

func getStatus() error {

	status, err := getAppStatus()
	if err != nil {
		return err
	}

	fmt.Printf("Heighliner application %s is ready!\n", status.AppName)
	fmt.Printf("You can access %s on %s [Username: %s Password: %s]\n\n", status.Cd.Provider, color.HiBlueString(status.Cd.URL),
		status.Cd.Username, status.Cd.Password)
	fmt.Printf("There are %d applications deployed by %s:\n", len(status.Apps), status.Cd.Provider)
	for i, info := range status.Apps {
		fmt.Printf("%d: %s\n", i+1, info.Name)
		if info.Service != nil {
			fmt.Printf("   %s has been deployed to k8s cluster, you can access it by k8s Service url: %s\n",
				info.Name, color.HiBlueString(info.Service.URL))
		}
		if info.Repo != nil {
			fmt.Printf("   %s's source code resides on %s repository: %s\n", info.Name, status.SCM.Provider, color.HiBlueString(info.Repo.URL))
		}
		if info.Username != "" && info.Password != "" {
			fmt.Printf("   credential: [Username: %s Password: %s]\n", info.Username, info.Password)
		}
		fmt.Println()
	}
	return nil
}
