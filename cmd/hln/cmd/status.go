package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/internal/k8sfactory"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/state/app"
)

func newStatusCmd(streams genericclioptions.IOStreams) *cobra.Command {
	c := &cobra.Command{
		Use:   "status [appName]",
		Short: "Show status of your application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := getStateInSpecificBackend()
			if err != nil {
				return err
			}

			apps, err := st.ListApps()
			if err != nil {
				return err
			}
			var found bool
			for _, s := range apps {
				if s == args[0] {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("application \"%s\" not found ", args[0])
			}
			return showStatus(streams.Out, args[0])
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

func showStatus(w io.Writer, appName string) error {

	status, err := getAppStatus(appName)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n🎉 Heighliner application %s is ready! ", status.AppName)

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
		fmt.Fprintf(w, "access URL: %s\n\n", color.HiBlueString(frontendService.Service.URL))
	}

	fmt.Fprintf(w, "There are %d services have been deployed:\n", len(status.UserServices))
	for _, info := range status.UserServices {
		fmt.Fprintf(w, "● %s\n", info.Service.Name)

		if info.Service.URL != "" {
			fmt.Fprintf(w, "  ● access URL: %s\n", color.HiBlueString(info.Service.URL))
		}

		if info.Repo != nil {
			fmt.Fprintf(w, "  ● resource code: %s\n", color.HiBlueString(info.Repo.URL))
		}

		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "There are %d addons have been deployed:\n", len(addonServices)+len(emptyAddonServices)+1)
	fmt.Fprintf(w, "● %s\n", status.CD.Provider)
	if status.CD.URL != "" {
		fmt.Fprintf(w, "  ● access URL: %s\n", color.HiBlueString(status.CD.URL))
	}
	if status.CD.Username != "" && status.CD.Password != "" {
		fmt.Fprintf(w, "  ● credential: [Username: %s Password: %s]\n", status.CD.Username, status.CD.Password)
	}
	fmt.Fprintln(w)

	for _, info := range addonServices {
		fmt.Fprintf(w, "● %s\n", info.Name)

		if info.URL != "" {
			fmt.Fprintf(w, "  ● access URL: %s\n", color.HiBlueString(info.URL))
		}

		if info.Username != "" && info.Password != "" {
			fmt.Fprintf(w, "  ● credential: [Username: %s Password: %s]\n", info.Username, info.Password)
		}

		if info.Prompt != "" {
			for _, prompt := range strings.Split(info.Prompt, ", ") {
				fmt.Fprintf(w, "  ● %s\n", prompt)
			}
		}

		fmt.Fprintln(w)
	}

	for _, info := range emptyAddonServices {
		fmt.Fprintf(w, "● %s\n", info.Name)
	}

	return nil
}
