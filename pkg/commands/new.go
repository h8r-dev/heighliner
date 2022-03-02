package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/state"
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/spf13/cobra"
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
	if val, ok := state.Stacks[projStack]; !ok {
		panic("no such stack")
	} else {
		err := initProj(projStack, "", val.Url)
		if err != nil {
			return err
		}
	}
	return nil
}

func initProj(name, dst, src string) error {
	if dst == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working dir: %w", err)
		}
		dst = cwd
	}
	// prepare hln dir
	ds, err := state.Make(dst)
	if err != nil {
		return err
	}
	s, err := ds.Find()
	if err == state.ErrNoStack {
		// prepare stack
		dir := path.Join(dst, "hln")
		s, err = state.New(name, dir, src)
		if err != nil {
			return err
		}
		err = s.Download()
		if err != nil {
			return err
		}
		err = s.Decompress()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// init project 1st step --dagger init
	err = util.Exec("dagger",
		"--project", dst,
		"init")
	if err != nil {
		return err
	}
	// 2nd step --dagger new hln -p /path/to/plans
	err = util.Exec("dagger",
		"--project", dst,
		"new", "hln",
		"-p", path.Join(s.Path, s.Name, "plans"))
	if err != nil {
		return err
	}
	return nil
}
