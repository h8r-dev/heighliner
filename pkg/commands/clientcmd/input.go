package clientcmd

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
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
	cmd := exec.Command(
		"dagger",
		"--project", "",
		"-e", "hln",
		"input", args[0], args[1], args[2])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
