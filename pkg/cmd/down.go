package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/project"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

func newDownCmd() *cobra.Command {
	var interactive bool

	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Take down your application",
		Args:  cobra.NoArgs,
	}

	downCmd.Flags().StringP("stack", "s", "", "Name of your stack")
	downCmd.Flags().StringArray("set", []string{}, "The input values of your project")
	downCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	downCmd.Flags().Bool("no-cache", false, "Disable caching")

	if err := downCmd.MarkFlagRequired("stack"); err != nil {
		log.Fatal().Err(err).Msg("Failed to mark flag required")
	}

	downCmd.Run = func(c *cobra.Command, args []string) {
		var err error
		lg := logger.New()

		stackName := c.Flags().Lookup("stack").Value.String()

		// Update the stack.
		s, err := stack.New(stackName)
		if err != nil {
			lg.Fatal().Err(err)
		}
		err = s.Update()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to update stack")
		}

		// Initialize the project.
		// Enter the project dir automatically.
		p := project.New(
			path.Join(state.GetCache(), path.Base(stackName)),
			path.Join(state.GetTemp(), path.Base(stackName)))
		err = p.Init()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to initialize project")
		}

		// Set input values.
		sch := schema.New()
		err = sch.AutomaticEnv(interactive)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to set automatic env")
		}

		// Execute the action.
		newArgs := []string{}
		newArgs = append(newArgs,
			"--log-format", viper.GetString("log-format"),
			"--log-level", viper.GetString("log-level"),
			"do", "down",
			"-p", "./plans")
		if c.Flags().Lookup("no-cache").Value.String() == "true" {
			newArgs = append(newArgs, "--no-cache")
		}
		err = util.Exec(
			util.Dagger,
			newArgs...)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute stack")
		}

		// Print the output.
		b, err := os.ReadFile("output.yaml")
		if err != nil {
			lg.Warn().Err(err).Msg("no output information")
		} else {
			fmt.Fprintf(os.Stdout, "\n%s", b)
		}
	}

	return downCmd
}
