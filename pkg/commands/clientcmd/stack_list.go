package clientcmd

import (
	"os"
	"os/exec"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	stackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List input values of the stack",
		Args:  cobra.NoArgs,
		RunE:  listStack,
	}
)

func listStack(c *cobra.Command, args []string) error {
	ds, err := datastore.Stat()
	if err != nil {
		return err
	}
	s, err := ds.Find()
	if err != nil {
		return err
	}
	cmd := exec.Command("dagger",
		"--project", s.Path,
		"-e", "hln",
		"input", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
