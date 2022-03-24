package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/hofstadter-io/hof/lib/mod"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a heighliner project",
		Args:  cobra.NoArgs,
		RunE:  newProj,
	}
)

func init() {
	newCmd.Flags().StringP("stack", "s", "", "The stack of your project")
	err := newCmd.MarkFlagRequired("stack")
	if err != nil {
		log.Fatal().Err(err)
	}
	if err := viper.BindPFlags(newCmd.Flags()); err != nil {
		panic(err)
	}
}

func newProj(c *cobra.Command, args []string) error {

	stackName := viper.GetString("stack")
	// Check if specified stack exist or not
	val, ok := stack.Stacks[stackName]
	if !ok {
		return stack.ErrNoSuchStack
	}

	s := stack.New(stackName)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current work dir: %w", err)
	}

	// Prepare stack files, Could be optimized with version function
	{
		err = os.RemoveAll(state.HeighlinerCacheHome)
		if err != nil {
			return err
		}

		err = s.Pull(val.URL, state.HeighlinerCacheHome)
		if err != nil {
			return err
		}
	}

	// TODO hide project dir
	// err := os.Chdir(filepath.Join(pwd, "project"))
	// if err != nil {
	// 	panic(err)
	// }
	err = copy.Copy(path.Join(state.HeighlinerCacheHome, stackName), pwd)
	if err != nil {
		return err
	}

	// $ hof mod vendor cue
	mod.InitLangs()
	err = mod.ProcessLangs("vendor", []string{"cue"})
	if err != nil {
		fmt.Println(err)
	}

	// Initialize & update project
	err = util.Exec("dagger", "project", "init")
	if err != nil {
		return err
	}
	err = util.Exec("dagger", "project", "update")
	if err != nil {
		return err
	}

	return nil
}
