package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// NewDownCmd creates and returns the down command of hln
func NewDownCmd() *cobra.Command {
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Shut down the application and clear resources",
		Args:  cobra.ArbitraryArgs,
		RunE:  downProj,
	}

	return downCmd
}

func downProj(c *cobra.Command, args []string) error {
	lg := logger.New()
	if err := state.EnterTemp(); err != nil {
		lg.Fatal().Err(err).Msg("Couldn't find project. Please run hln new to create one")
	}
	return util.Exec("dagger",
		"--log-format", viper.GetString("log-format"),
		"--log-level", viper.GetString("log-level"),
		"-p", "./plans",
		"do", "down")
}
