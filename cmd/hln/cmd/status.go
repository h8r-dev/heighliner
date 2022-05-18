package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/internal/k8sfactory"
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
	kubecli, err := k8sfactory.GetDefaultClientSet()
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

	fmt.Printf("Heighliner application %s is ready!", status.AppName)

	var frontendService *app.ServiceInfo
	var addonServices []app.ServiceInfo
	var userServices []app.ServiceInfo
	var emptyAddonServices []app.ServiceInfo
	for _, info := range status.Services {
		if info.Type == "frontend" && info.Service != nil && info.Service.URL != "" {
			frontendService = &info
			userServices = append(userServices, info)
			continue
		}

		if info.Infra == "true" {
			if info.URL == "" && info.Service == nil && info.Repo == nil {
				emptyAddonServices = append(emptyAddonServices, info)
			} else {
				addonServices = append(addonServices, info)
			}
			continue
		}

		userServices = append(userServices, info)
	}

	if frontendService != nil {
		fmt.Printf("access URL: %s", frontendService.Service.URL)
	}
	fmt.Println()

	//fmt.Printf("You can access %s on %s [Username: %s Password: %s]\n\n", status.CD.Provider, color.HiBlueString(status.CD.URL),
	//	status.CD.Username, status.CD.Password)
	fmt.Printf("There are %d services have been deployed:\n", len(status.Services))
	for _, info := range userServices {
		fmt.Printf("● %s\n", info.Name)
		if info.URL != "" {
			fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.URL))
		} else if info.Service != nil {
			fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.Service.URL))
		}
		//if info.Service != nil {
		//	fmt.Printf("   %s has been deployed to k8s cluster, you can access it by k8s Service url: %s\n",
		//		info.Name, color.HiBlueString(info.Service.URL))
		//}
		if info.Repo != nil {
			fmt.Printf("  ● resource code: %s\n", color.HiBlueString(info.Repo.URL))
		}
		if info.Username != "" && info.Password != "" {
			fmt.Printf("  ● credential: [Username: %s Password: %s]\n", info.Username, info.Password)
		}
		//if info.Prompt != "" {
		//	fmt.Printf("   %s\n", info.Prompt)
		//}
		fmt.Println()
	}

	if len(addonServices)+len(emptyAddonServices) > 0 {

		fmt.Printf("There are %d addons have been deployed:\n", len(addonServices)+len(emptyAddonServices))
		for _, info := range addonServices {
			fmt.Printf("● %s\n", info.Name)
			if info.URL != "" {
				fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.URL))
			} else if info.Service != nil {
				fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.Service.URL))
			}
			//if info.Service != nil {
			//	fmt.Printf("   %s has been deployed to k8s cluster, you can access it by k8s Service url: %s\n",
			//		info.Name, color.HiBlueString(info.Service.URL))
			//}
			if info.Repo != nil {
				fmt.Printf("  ● resource code: %s\n", color.HiBlueString(info.Repo.URL))
			}
			if info.Username != "" && info.Password != "" {
				fmt.Printf("  ● credential: [Username: %s Password: %s]\n", info.Username, info.Password)
			}
			if info.Prompt != "" {
				fmt.Printf("  ● %s\n", info.Prompt)
			}
			fmt.Println()
		}
		for _, info := range emptyAddonServices {
			fmt.Printf("● %s\n", info.Name)
		}
	}

	return nil
}
