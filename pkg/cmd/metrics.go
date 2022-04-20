package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Metrics represents the monitoring metrics.
type Metrics struct {
	Infras []Infra `yaml:"infra"`
}

// Infra represents a component of the infrastructure.
type Infra struct {
	Type     string `yaml:"type"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func newMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show dashboard of monitoring metrics",
	}

	cmd.RunE = func(c *cobra.Command, args []string) error {
		printTarget := os.Stdout
		b, err := os.ReadFile(appInfo)
		if err != nil {
			return err
		}
		m := new(Metrics)
		if err := yaml.Unmarshal(b, m); err != nil {
			return err
		}
		for _, infra := range m.Infras {
			if infra.Type == "grafana" {
				u := url.URL{
					Scheme:   "http",
					Host:     infra.URL,
					Path:     "explore",
					RawQuery: `left={"datasource"="Loki"}`,
				}
				fmt.Fprintf(printTarget, "URL: %s\nUsername: %s\nPassword: %s\n", u.String(), infra.Username, infra.Password)
				if err := browser.OpenURL(u.String()); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return cmd
}
