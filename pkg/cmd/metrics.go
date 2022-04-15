package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/h8r-dev/heighliner/pkg/logger"
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

	cmd.Run = func(c *cobra.Command, args []string) {
		lg := logger.New()
		printTarget := os.Stdout
		b, err := os.ReadFile(appInfo)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to read app state")
		}
		m := new(Metrics)
		err = yaml.Unmarshal(b, m)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to marsahl app state")
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
				err := browser.OpenURL(u.String())
				if err != nil {
					lg.Fatal().Err(err).Msg("failed to open browser")
				}
			}
		}
	}

	return cmd
}
