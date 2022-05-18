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

	fmt.Printf("🎉 Heighliner application %s is ready! ", status.AppName)

	var frontendService app.UserService
	var addonServices []app.ServiceInfo
	var emptyAddonServices []app.ServiceInfo
	for _, info := range status.Services {

		if info.Infra == "true" {
			if info.URL == "" {
				emptyAddonServices = append(emptyAddonServices, info)
			} else {
				addonServices = append(addonServices, info)
			}
			continue
		}
	}

	var found bool
	for _, service := range status.UserServices {
		if service.Type == "frontend" {
			frontendService = service
			found = true
			break
		}
	}

	if found {
		fmt.Printf("access URL: %s", color.HiBlueString(frontendService.Service.URL))
	}
	fmt.Printf("\n\n")

	fmt.Printf("There are %d services have been deployed:\n", len(status.UserServices))
	for _, info := range status.UserServices {
		fmt.Printf("● %s\n", info.Service.Name)

		if info.Service.URL != "" {
			fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.Service.URL))
		}

		if info.Repo != nil {
			fmt.Printf("  ● resource code: %s\n", color.HiBlueString(info.Repo.URL))
		}

		fmt.Println()
	}

	fmt.Printf("There are %d addons have been deployed:\n", len(addonServices)+len(emptyAddonServices)+1)
	fmt.Printf("● %s\n", status.CD.Provider)
	if status.CD.URL != "" {
		fmt.Printf("  ● access URL: %s\n", color.HiBlueString(status.CD.URL))
	}
	if status.CD.Username != "" && status.CD.Password != "" {
		fmt.Printf("  ● credential: [Username: %s Password: %s]\n", status.CD.Username, status.CD.Password)
	}
	fmt.Println()

	for _, info := range addonServices {
		fmt.Printf("● %s\n", info.Name)

		if info.URL != "" {
			fmt.Printf("  ● access URL: %s\n", color.HiBlueString(info.URL))
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

	return nil
}
