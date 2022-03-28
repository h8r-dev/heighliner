package cmd

import (
	"path"

	"github.com/hofstadter-io/hof/lib/mod"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
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

	if err := state.InitTemp(); err != nil {
		lg.Fatal().Err(err).Msg("There is already a project. Consider drop it?")
	}

	// Check if target stack exists or not
	stackName := viper.GetString("stack")
	if stackName == "" {
		state.CleanTemp()
		lg.Fatal().Msg("Please specify a stack with -s flag")
	}
	s, err := stack.New(stackName)
	if err != nil {
		state.CleanTemp()
		lg.Fatal().Err(err).Msgf("failed to create project with stack \"%s\"", stackName)
	}

	// Fetch target stack
	if err := state.CleanCache(); err != nil {
		lg.Fatal().Err(err).Msg("failed to clean stack cache")
	}
	if err := s.Pull(); err != nil {
		lg.Fatal().Err(err).Msgf("failed to pull stack %s", s.Name)
	}
	if err := s.Copy(path.Join(state.Cache, s.Name), state.Temp); err != nil {
		lg.Fatal().Err(err).Msg("failed to copy stack")
	}

	// Change diretory
	if err := state.EnterTemp(); err != nil {
		lg.Fatal().Err(err).Msg("failed to enter project dir")
	}

	// $ hof mod vendor cue
	mod.InitLangs()
	if err := mod.ProcessLangs("vendor", []string{"cue"}); err != nil {
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
