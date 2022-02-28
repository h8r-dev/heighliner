package clientcmd

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/rs/zerolog/log"
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
		panic(err)
	}
}

func newProj(c *cobra.Command, args []string) error {
	if val, ok := stack.Stacks[projStack]; !ok {
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
			return errors.New("failed to get current working dir")
		}
		dst = cwd
	}
	// prepare hln dir
	ds, err := datastore.Make(dst)
	if err != nil {
		return err
	}
	s, err := ds.Find()
	if err == datastore.ErrNoStack {
		// prepare stack
		dir := path.Join(dst, "hln")
		s, err = stack.New(name, dir, src)
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
	cmd := exec.Command("dagger",
		"--project", dst,
		"init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	// 2nd step --dagger new hln -p /path/to/plans
	cmd = exec.Command("dagger",
		"--project", dst,
		"new", "hln",
		"-p", path.Join(s.Path, s.Name, "plans"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
