package clientcmd

import (
	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Manage secrets",
		Args:  cobra.NoArgs,
		RunE:  initDataStore,
	}
)

func initDataStore(cmd *cobra.Command, args []string) error {
	_, err := datastore.Init()
	return err
}
