package cmd

import (
	"fmt"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/pkg/errors"

	"github.com/fatih/color"
	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "status [appName]",
		Short: "Show status of your application",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return getStatus(args[0])
		},
	}

	return c
}

// GetTfProvider For hln down
func GetTfProvider(appName string) (string, error) {
	cs, err := getConfigMapState()
	if err != nil {
		return "", err
	}

	return cs.LoadTfProvider(appName)
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

	cs, err := getConfigMapState()
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

func getStatus(appName string) error {

	status, err := getAppStatus(appName)
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
