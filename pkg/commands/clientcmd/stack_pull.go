package clientcmd

import (
	"errors"

	"github.com/h8r-dev/heighliner/pkg/commands/clientcmd/state"

	"github.com/spf13/cobra"
)

var (
	stackPullCmd = &cobra.Command{
		Use:   "pull [stack name] [stack url]",
		Short: "Pull a stack",
		Long:  "",
		RunE:  pullStack,
	}
)

func pullStack(c *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("please specify stack name and url")
	}

	hs, err := state.InitHlnStore()
	if err != nil {
		return err
	}

	stack, err := hs.NewStack(args[0], args[1])
	if err != nil {
		return err
	}

	err = stack.Download()
	if err != nil {
		return err
	}

	err = stack.Decompress()
	if err != nil {
		return err
	}

	return nil
}
