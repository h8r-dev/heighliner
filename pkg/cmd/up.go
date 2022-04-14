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
	"github.com/h8r-dev/heighliner/pkg/util/dagger"
)

func newUpCmd() *cobra.Command {
	var interactive bool

	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up your application",
		Args:  cobra.NoArgs,
	}

	upCmd.Flags().StringP("stack", "s", "", "Name of your stack")
	upCmd.Flags().String("dir", "", "Path to your local stack")
	upCmd.Flags().StringArray("set", []string{}, "The input values of your project")
	upCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	upCmd.Flags().Bool("no-cache", false, "Disable caching")

	upCmd.Run = func(c *cobra.Command, args []string) {
		var (
			err error
			p   *project.Project
		)
		lg := logger.New()

		// Validate args
		stackName := c.Flags().Lookup("stack").Value.String()
		dir := c.Flags().Lookup("dir").Value.String()
		switch {
		case dir != "" && stackName != "":
			lg.Fatal().Msg("please do not specify both stack and dir at the same time")
		case stackName != "":
			// Update the stack.
			s, err := stack.New(stackName)
			if err != nil {
				lg.Fatal().Err(err).Msg("no such stack")
			}
			err = s.Update()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to update stack")
			}
			p = project.New(
				path.Join(state.GetCache(), path.Base(stackName)),
				path.Join(state.GetTemp(), path.Base(stackName)))
		default:
			stackSrc := util.Abs(dir)
			lg.Info().Msgf("Using local stack %s", stackSrc)
			fi, err := os.Stat(stackSrc)
			if err != nil {
				lg.Fatal().Err(err).Msgf("failed to find stack in %s", stackSrc)
			}
			if !fi.IsDir() {
				lg.Fatal().Msgf("%s is not a directory", stackSrc)
			}
			p = project.New(stackSrc, path.Join(state.GetTemp(), path.Base(stackSrc)))
		}

		// Initialize the project.
		// Enter the project dir automatically.

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
			"do", "up",
			"-p", "./plans")
		if c.Flags().Lookup("no-cache").Value.String() == "true" {
			newArgs = append(newArgs, "--no-cache")
		}
		err = util.Exec(
			dagger.GetPath(),
			newArgs...)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to execute stack")
		}

		// Print the output.
		b, err := os.ReadFile(appOutputPath)
		if err != nil {
			lg.Warn().Err(err).Msg("no output information")
		} else {
			fmt.Fprintf(os.Stdout, "\n%s", b)
		}
	}

	return upCmd
}
