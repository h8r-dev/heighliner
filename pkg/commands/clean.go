package commands

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/state"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "clean all caches",
		Args:  cobra.NoArgs,
		RunE:  cleanHeighlinerCache,
	}
)

func cleanHeighlinerCache(c *cobra.Command, args []string) error {
	err := state.CleanHeighlinerCaches()
	if err != nil {
		return err
	}
	return nil
}
