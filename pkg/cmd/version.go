package cmd

import (
	"fmt"
	"runtime"

	"github.com/h8r-dev/heighliner/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Drop the current heighliner project",
		Args:  cobra.NoArgs,
		Run:   printVersion,
	}
)

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("hln %s (%s) %s/%s\n",
		version.Version,
		version.Revision,
		runtime.GOOS, runtime.GOARCH,
	)
}
