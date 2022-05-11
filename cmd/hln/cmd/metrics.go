package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type metricsOptions struct {
	genericclioptions.IOStreams
}

type Metrics struct {
	AppName       string
	CridentialRef Cridential
	DashboardRefs []MonitorDashboard
}

type Cridential struct {
	Username string
	Password string
}

type MonitorDashboard struct {
	Title string
	URL   url.URL
}

func (o *metricsOptions) Run(args []string) error {
	appName := args[0]
	metrics, err := getMetrics(appName)
	if err != nil {
		return fmt.Errorf("failed to get application metrics: %w", err)
	}
	showMetrics(o.Out, metrics)
	return nil
}

func newMetricsCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &metricsOptions{
		IOStreams: streams,
	}

	cmd := &cobra.Command{
		Use:   "metrics [appName]",
		Short: "Show dashboard of monitoring metrics",
		Args:  cobra.ExactArgs(1),
	}

	cmd.RunE = func(c *cobra.Command, args []string) error {
		return o.Run(args)
	}

	return cmd
}

func getMetrics(appName string) (*Metrics, error) {
	st, err := getStateInSpecificBackend()
	if err != nil {
		return nil, err
	}
	ao, err := st.LoadOutput(appName)
	if err != nil {
		return nil, err
	}
	metrics := &Metrics{
		AppName: ao.ApplicationRef.Name,
	}
	var foundFlag bool // false by default
	for _, argoApp := range ao.CD.ApplicationRef {
		if argoApp.Type == "monitoring" {
			foundFlag = true
			metrics.CridentialRef.Username = argoApp.Username
			metrics.CridentialRef.Password = argoApp.Password
			if argoApp.Annotations != "" {
				type MDashboard struct {
					Title string
					Path  string
				}
				mdb := MDashboard{}
				if err := json.Unmarshal([]byte(argoApp.Annotations), &mdb); err != nil {
					return nil, err
				}
				metrics.DashboardRefs = append(metrics.DashboardRefs, MonitorDashboard{
					Title: mdb.Title,
					URL: url.URL{
						Host: argoApp.URL,
						Path: mdb.Path,
					},
				})
			}
		}
	}
	if !foundFlag {
		return nil, errors.New("target app doesn't have any monitor component")
	}
	return metrics, nil
}

func showMetrics(w io.Writer, m *Metrics) {
	fmt.Fprintf(w, "The metrics of %s:\n", m.AppName)
	fmt.Fprintf(w, "Cridentials for login:\n")
	fmt.Fprintf(w, "  Username: %s\n", color.HiBlueString(m.CridentialRef.Username))
	fmt.Fprintf(w, "  Password: %s\n", color.HiBlueString(m.CridentialRef.Password))
	for _, db := range m.DashboardRefs {
		fmt.Fprintf(w, "Dashboard %s:\n", db.Title)
		fmt.Fprintf(w, "  URL: %s\n", color.CyanString(db.URL.String()))
	}
}
