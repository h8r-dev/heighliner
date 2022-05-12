package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/state/app"
)

func newStatusCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "status [appName]",
		Short: "Show status of your application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showStatus(args[0])
		},
	}

	return c
}

// GetTFProvider For hln down
func GetTFProvider(appName string) (string, error) {
	cs, err := getStateInSpecificBackend()
	if err != nil {
		return "", err
	}

	return cs.LoadTFProvider(appName)
}

// Get state in specific backend by env, such as: CONFIG_MAP, S3, LOCAL_FILE
func getStateInSpecificBackend() (state.State, error) {
	if l, ok := os.LookupEnv("STATE_BACKEND"); ok && l == "LOCAL_FILE" {
		return &state.LocalFileState{}, nil
	}
	return getConfigMapState()
}

func getConfigMapState() (state.State, error) {
	kubecli, err := getDefaultClientSet()
	if err != nil {
		return nil, fmt.Errorf("failed to make kube client: %w", err)
	}

	cs := state.ConfigMapState{ClientSet: kubecli}
	return &cs, nil
}

// Get Heighliner application status from k8s configmap
func getAppStatus(appName string) (*app.Status, error) {

	cs, err := getStateInSpecificBackend()
	if err != nil {
		return nil, err
	}

	ao, err := cs.LoadOutput(appName)
	if err != nil {
		return nil, err
	}

	s := ao.ConvertOutputToStatus()
	s.AppName = appName
	return &s, nil
}

func showStatus(appName string) error {

	status, err := getAppStatus(appName)
	if err != nil {
		return err
	}

	fmt.Printf("Heighliner application %s is ready!\n", status.AppName)
	fmt.Printf("You can access %s on %s [Username: %s Password: %s]\n\n", status.CD.Provider, color.HiBlueString(status.CD.URL),
		status.CD.Username, status.CD.Password)
	fmt.Printf("There are %d services deployed by %s:\n", len(status.Services), status.CD.Provider)
	for i, info := range status.Services {
		fmt.Printf("%d: %s\n", i+1, info.Name)
		if info.URL != "" {
			fmt.Printf("   You can access %s from broswer by url: %s\n",
				info.Name, color.HiBlueString(info.URL))
		}
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
		if info.Prompt != "" {
			fmt.Printf("   %s\n", info.Prompt)
		}
		fmt.Println()
	}
	return nil
}
