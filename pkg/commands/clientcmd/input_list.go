package clientcmd

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
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
	cmd := exec.Command(
		"dagger",
		"--project", "",
		"-e", "hln",
		"input", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
