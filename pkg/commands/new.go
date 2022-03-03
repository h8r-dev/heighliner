package commands

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/state"
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
)

var (
	projStack string

	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a heighliner project",
		Args:  cobra.NoArgs,
		RunE:  newProj,
	}
)

func init() {
	newCmd.Flags().StringVarP(&projStack, "stack", "s", "", "The stack of your project")
	err := newCmd.MarkFlagRequired("stack")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func newProj(c *cobra.Command, args []string) error {
	val, ok := state.Stacks[projStack]
	if !ok {
		return fmt.Errorf("no such stack")
	}
	err := initProj(projStack, "", val.URL)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully create new project with stack: %s\n", projStack)
	return nil
}

func initProj(name, dest, src string) error {
	if dest == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working dir: %w", err)
		}
		dest = cwd
	}
	s := state.NewStack(name)
	err := s.Check()
	if err != nil {
		if errors.Is(err, state.ErrStackNotExist) {
			err := s.Pull(src)
			if err != nil {
				return fmt.Errorf("failed to pull stack %s: %w", s.Name, err)
			}
		} else {
			return err
		}
	}
	err = s.Copy(dest)
	if err != nil {
		return fmt.Errorf("failed to copy stack %s: %w", s.Name, err)
	}
	// init project 1st step --dagger init
	err = util.Exec("dagger",
		"--project", dest,
		"init")
	if err != nil {
		return err
	}
	// 2nd step --dagger new hln -p /path/to/plans
	err = util.Exec("dagger",
		"--project", dest,
		"new", "hln",
		"-p", path.Join(dest, "plans"))
	if err != nil {
		return err
	}
	return nil
}
