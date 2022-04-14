package cmd

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/project"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/dagger"
)

func newTestCmd() *cobra.Command {
	var interactive bool

	cmd := &cobra.Command{
		Use:    "test",
		Short:  "Test your stack",
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	cmd.Flags().String("dir", "", "Path to your local stack")
	cmd.Flags().StringP("plan", "p", "", "Relative path to your test plan")
	cmd.Flags().StringArray("set", []string{}, "The input values of your project")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	cmd.Flags().Bool("no-cache", false, "Disable caching")

	cmd.Run = func(c *cobra.Command, args []string) {
		var err error
		lg := logger.New()

		dir := c.Flags().Lookup("dir").Value.String()

		// If stack flag is not set, use the current directory as stack source.
		if dir == "" {
			dir, err = os.Getwd()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to get current working directory")
			}
		}

		// Create a project object.
		dir, err = homedir.Expand(dir)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to expand path")
		}
		lg.Info().Msgf("Using local stack %s", dir)
		fi, err := os.Stat(dir)
		if err != nil {
			lg.Fatal().Err(err).Msgf("failed to find stack in %s", dir)
		}
		if !fi.IsDir() {
			lg.Fatal().Msgf("%s is not a directory", dir)
		}
		p := project.New(dir, path.Join(state.GetTemp(), path.Base(dir)))

		// Initialize the project.
		// Enter the project dir automatically.
		err = p.Init()
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to initialize project")
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
		err = util.Exec(
			dagger.GetPath(),
			newArgs...)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute stack")
		}
	}

	return cmd
}
