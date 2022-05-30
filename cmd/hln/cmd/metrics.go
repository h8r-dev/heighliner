package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"
)

type metricsOptions struct {
	genericclioptions.IOStreams
}

// Metrics to print
type Metrics struct {
	AppName       string
	CredentialRef Credential
	DashboardRefs []*MonitorDashboard
}

// Credential for login
type Credential struct {
	Username string
	Password string
}

// MonitorDashboard of apps
type MonitorDashboard struct {
	Title string
	URL   string
}

func (o *metricsOptions) Run(appName string) error {
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
		return o.Run(args[0])
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
			metrics.CredentialRef.Username = argoApp.Username
			metrics.CredentialRef.Password = argoApp.Password
			if argoApp.Annotations != "" {
				str := argoApp.Annotations
				data, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					return nil, fmt.Errorf("failed to decode annotations :%w", err)
				}
				mdbs := []struct {
					Title string `json:"title"`
					Path  string `json:"path"`
				}{}
				if err := yaml.Unmarshal(data, &mdbs); err != nil {
					return nil, fmt.Errorf("bad annotations format: %w", err)
				}
				for _, mdb := range mdbs {
					metrics.DashboardRefs = append(metrics.DashboardRefs, &MonitorDashboard{
						Title: mdb.Title,
						URL:   argoApp.URL + mdb.Path,
					})
				}
			}
		}
	}
	if !foundFlag {
		return nil, errors.New("target app doesn't have any monitor component")
	}
	return metrics, nil
}

func showMetrics(w io.Writer, m *Metrics) {
	fmt.Fprintf(w, "Use this credential to login the monitoring dashboards of %s:\n", m.AppName)
	fmt.Fprintf(w, "  Username: %s\n", color.HiBlueString(m.CredentialRef.Username))
	fmt.Fprintf(w, "  Password: %s\n", color.HiBlueString(m.CredentialRef.Password))
	fmt.Fprintf(w, "\nApplication %s has %d available dashboard(s):\n", m.AppName, len(m.DashboardRefs))
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	defer func() {
		err := tw.Flush()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()
	fmt.Fprintf(tw, "NAME\tURL\n")
	for _, db := range m.DashboardRefs {
		fmt.Fprintf(tw, "%s\t%s\n", db.Title, color.CyanString(db.URL))
	}
}
