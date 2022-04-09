package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/project"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

func newUpCmd() *cobra.Command {
	var interactive bool

	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up your application",
		Args:  cobra.NoArgs,
	}

	upCmd.Flags().StringP("stack", "s", "", "Stack name or path")
	upCmd.Flags().StringArray("set", []string{}, "The input values of your project")
	upCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	upCmd.Flags().Bool("no-cache", false, "Disable caching")

	upCmd.RunE = func(c *cobra.Command, args []string) error {
		var err error
		lg := logger.New()

		stackSrc := c.Flags().Lookup("stack").Value.String()

		// Create a project object.
		if _, ok := stack.Stacks[stackSrc]; ok {
			lg.Info().Msgf("Using official stack %s", stackSrc)
			s, err := stack.New(stackSrc)
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to create stack")
			}
			err = s.Pull()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to pull stack")
			}
			stackSrc = path.Join(state.GetCache(), stackSrc)
		} else {
			lg.Fatal().Msgf("can not find stack %s", stackSrc)
		}
		p := project.New(stackSrc, path.Join(state.GetTemp(), path.Base(stackSrc)))

		// Initialize the project.
		// Enter the project dir automatically.
		state.CleanTemp()
		err = p.Init()
		defer p.Clean()
		if err != nil {
			p.Clean()
			lg.Fatal().Err(err).Msg("failed to initialize project")
		}

		// Set input values.
		sch := schema.New()
		err = sch.AutomaticEnv(interactive)
		if err != nil {
			return fmt.Errorf("failed to set env: %w", err)
		}

		// Execute the action.
		newArgs := []string{}
		newArgs = append(newArgs,
			"--log-format", viper.GetString("log-format"),
			"--log-level", viper.GetString("log-level"),
			"-p", "./plans",
			"do", "up")
		if c.Flags().Lookup("no-cache").Value.String() == "true" {
			newArgs = append(newArgs, "--no-cache")
		}
		if err := util.Exec("dagger", newArgs...); err != nil {
			return err
		}

		// Print the output.
		b, err := os.ReadFile("output.yaml")
		if err != nil {
			lg.Warn().Err(err).Msg("no output information")
		} else {
			fmt.Fprintf(os.Stdout, "\n%s", b)
		}
		return nil
	}

	return upCmd
}
