package clientcmd

import (
	"os"
	"os/exec"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	stackInputCmd = &cobra.Command{
		Use:   "input",
		Short: "input a value for a stack",
		Args:  cobra.ExactArgs(3),
		RunE:  stackInput,
	}
)

func stackInput(c *cobra.Command, args []string) error {
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
		"-e", "hln",
		"input", args[0], args[1], args[2])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
