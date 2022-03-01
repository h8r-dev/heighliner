package commands

import (
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/spf13/cobra"
)

var (
	inputListCmd = &cobra.Command{
		Use:   "list",
		Short: "List input values",
		Args:  cobra.NoArgs,
		RunE:  inputList,
	}
)

func inputList(c *cobra.Command, args []string) error {
	err := util.Exec(
		"dagger",
		"--project", "",
		"-e", "hln",
		"input", "list")
	if err != nil {
		return err
	}
	return nil
}
