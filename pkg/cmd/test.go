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
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// ActionGenerator generates a command with a description and a callback.
func newTestCmd() *cobra.Command {
	var interactive bool

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test your stack",
		Args:  cobra.NoArgs,
	}

	testCmd.Flags().StringP("stack", "s", "", "Path to your stack directory")
	testCmd.Flags().StringP("plan", "p", "", "Relative path to your test plan")
	testCmd.Flags().StringArray("set", []string{}, "The input values of your project")
	testCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	testCmd.Flags().Bool("no-cache", false, "Disable caching")

	testCmd.RunE = func(c *cobra.Command, args []string) error {
		var err error
		lg := logger.New()

		stackSrc := c.Flags().Lookup("stack").Value.String()

		// If stack flag is not set, use the current directory as stack source.
		if stackSrc == "" {
			stackSrc, err = os.Getwd()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to get current working directory")
			}
		}

		// Create a project object.
		stackSrc = util.Abs(stackSrc)
		lg.Info().Msgf("Using local stack %s", stackSrc)
		fi, err := os.Stat(stackSrc)
		if err != nil {
			lg.Fatal().Err(err).Msgf("failed to find stack in %s", stackSrc)
		}
		if !fi.IsDir() {
			lg.Fatal().Msgf("%s is not a directory", stackSrc)
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
		interactive := c.Flags().Lookup("interactive").Value.String()
		err = sch.AutomaticEnv(interactive == "true")
		if err != nil {
			return fmt.Errorf("failed to set env: %w", err)
		}

		// Execute the action.
		newArgs := []string{}
		newArgs = append(newArgs,
			"--log-format", viper.GetString("log-format"),
			"--log-level", viper.GetString("log-level"),
			"do", "test")
		if plan := c.Flags().Lookup("plan").Value.String(); plan != "" {
			newArgs = append(newArgs, "--plan", plan)
		}
		if c.Flags().Lookup("no-cache").Value.String() == "true" {
			newArgs = append(newArgs, "--no-cache")
		}
		if err := util.Exec("dagger", newArgs...); err != nil {
			return err
		}

		// Print the output.
		b, err := os.ReadFile("output.yaml")
		if err != nil {
			lg.Info().Err(err).Msg("no output information")
		} else {
			fmt.Fprintf(os.Stdout, "\n%s", b)
		}
		return nil
	}

	return testCmd
}
