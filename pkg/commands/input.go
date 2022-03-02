package commands

import (
	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/spf13/cobra"
)

var (
	inputCmd = &cobra.Command{
		Use:   "input [type] [name] [value]",
		Short: "Input a value",
		RunE:  inputValue,
	}
)

func init() {
	inputCmd.AddCommand(
		inputListCmd,
	)
}

func inputValue(c *cobra.Command, args []string) error {
	err := util.Exec(
		"dagger",
		"--project", "",
		"-e", "hln",
		"input", args[0], args[1], args[2])
	if err != nil {
		return err
	}
	return nil
}
