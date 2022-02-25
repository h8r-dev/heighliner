package clientcmd

import (
	"os"
	"os/exec"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		RunE:  upStack,
	}
)

func upStack(c *cobra.Command, args []string) error {
	ds, err := datastore.Stat()
	if err != nil {
		return err
	}
	s, err := ds.Find()
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"dagger",
		"--project", s.Path,
		"up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
