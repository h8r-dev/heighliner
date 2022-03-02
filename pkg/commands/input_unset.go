package commands

import (
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/spf13/cobra"
)

var (
	inputUnsetCmd = &cobra.Command{
		Use:   "unset [name]",
		Short: "unset an input value",
		Args:  cobra.ExactArgs(1),
		RunE:  inputUnset,
	}
)

func inputUnset(c *cobra.Command, args []string) error {
	err := util.Exec(
		"dagger",
		"--project", "",
		"-e", "hln",
		"input", "unset",
		args[0],
	)
	if err != nil {
		return err
	}
	return nil
}
