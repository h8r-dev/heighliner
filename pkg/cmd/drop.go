package cmd

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/state"
)

var (
	dropCmd = &cobra.Command{
		Use:   "drop",
		Short: "Drop the current heighliner project",
		Args:  cobra.NoArgs,
		Run:   dropProj,
	}
)

func dropProj(cmd *cobra.Command, args []string) {
	state.CleanTemp()
}
