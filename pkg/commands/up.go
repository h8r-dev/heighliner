package commands

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
)

var (
	logFormat string
	upCmd     = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		Args:  cobra.NoArgs,
		RunE:  upStack,
	}
)

func init() {
	upCmd.Flags().StringVarP(&logFormat, "log-format", "", "auto", `Log format (auto, plain, tty, json) (default "auto")`)
}

func upStack(c *cobra.Command, args []string) error {
	err := util.Exec(
		"dagger",
		"--project", "",
		"--log-format", logFormat,
		"up",
	)
	return err
}
