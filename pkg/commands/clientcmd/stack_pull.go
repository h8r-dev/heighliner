package clientcmd

import (
	"errors"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/spf13/cobra"
)

var (
	stackPullCmd = &cobra.Command{
		Use:   "pull [stack name] (stack url)",
		Short: "Pull a stack",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  pullStack,
	}
)

func pullStack(c *cobra.Command, args []string) error {

	ds, err := datastore.Stat()
	if err != nil {
		return err
	}
	var s *stack.Stack
	if len(args) == 1 {
		switch args[0] {
		case "sample":
			s, err = stack.New(
				args[0],
				ds.Path,
				"https://stack.h8r.io/sample-latest.tar.gz")
		default:
			return errors.New("can not find such stack")
		}
	} else {
		s, err = stack.New(args[0], ds.Path, args[1])
	}
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

	return nil
}
