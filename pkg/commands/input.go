package commands

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
)

var (
	srcFile string

	inputCmd = &cobra.Command{
		Use:   "input [type] [name] [value]",
		Short: "Input a value",
		Args:  cobra.ArbitraryArgs,
		RunE:  inputValue,
	}
)

func init() {
	inputCmd.Flags().StringVarP(&srcFile, "file", "f", "", "Source file to read from")
	inputCmd.AddCommand(
		inputListCmd,
		inputUnsetCmd,
	)
}

func inputValue(c *cobra.Command, args []string) error {
	var err error
	switch {
	case len(args) == 2:
		if srcFile == "" {
			return errors.New("please specify input source")
		}
		err = util.Exec(
			"dagger",
			"--project", "",
			"-e", "hln",
			"input",
			args[0], args[1],
			"-f", srcFile,
		)
	case len(args) == 3:
		err = util.Exec(
			"dagger",
			"--project", "",
			"-e", "hln",
			"input",
			args[0], args[1], args[2],
		)
	default:
		err = util.Exec(
			"dagger",
			"input",
		)
	}
	if err != nil {
		return err
	}
	return nil
}
