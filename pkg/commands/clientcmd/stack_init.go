package clientcmd

import (
	"os"
	"os/exec"
	"path"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	stackInitCmd = &cobra.Command{
		Use:   "init",
		Short: "init a dagger plan",
		Args:  cobra.NoArgs,
		RunE:  initStack,
	}
)

func initStack(c *cobra.Command, args []string) error {
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
		"init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal().Msg(string(out))
	}
	cmd = exec.Command("dagger",
		"--project", s.Path,
		"new", "hln",
		"-p", path.Join(s.Path, "plans"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
