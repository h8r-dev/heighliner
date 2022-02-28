package clientcmd

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		Args:  cobra.NoArgs,
		RunE:  upStack,
	}
)

func upStack(c *cobra.Command, args []string) error {
	cmd := exec.Command(
		"dagger",
		"--project", "",
		"up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
