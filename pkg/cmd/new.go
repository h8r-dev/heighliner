package cmd

import (
	"github.com/hofstadter-io/hof/lib/mod"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/proj"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a heighliner project",
		Args:  cobra.NoArgs,
		Run:   newProj,
	}
)

func init() {
	newCmd.Flags().StringP("stack", "s", "", "The stack of your project")
	if err := viper.BindPFlags(newCmd.Flags()); err != nil {
		log.Fatal().Err(err).Msg("failed to bind flags")
	}
}

func newProj(c *cobra.Command, args []string) {
	var lg = logger.New()

	stackName := viper.GetString("stack")
	if stackName == "" {
		lg.Fatal().Msg("Please specify a stack with -s flag")
	}

	s, err := stack.New(stackName)
	if err != nil {
		lg.Fatal().Err(err).Msgf("failed to create project with stack %s", stackName)
	}

	p := proj.New(s, state.NewTemp())
	if err := p.Init(); err != nil {
		lg.Fatal().Err(err).Msgf("failed to initialize project")
	}

	// $ hof mod vendor cue
	mod.InitLangs()
	err = mod.ProcessLangs("vendor", []string{"cue"})
	if err != nil {
		lg.Warn().Err(err).Msg("failed to fetch cuemods")
	}

	// Initialize & update project
	err = util.Exec("dagger", "project", "init")
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to execute dagger command")
	}
	err = util.Exec("dagger", "project", "update")
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to execute dagger command")
	}
}
