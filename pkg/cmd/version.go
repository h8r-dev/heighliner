package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/version"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
	}

	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("hln %s (%s) %s/%s\n",
			version.Version,
			version.Revision,
			runtime.GOOS, runtime.GOARCH,
		)
	}

	return versionCmd
}
