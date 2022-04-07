package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/version"
)

// NewVersionCmd creates and returns the version command of hln
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		Run:   printVersion,
	}

	return versionCmd
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("hln %s (%s) %s/%s\n",
		version.Version,
		version.Revision,
		runtime.GOOS, runtime.GOARCH,
	)
}
