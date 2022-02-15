package clientcmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	stackPullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull a stack",
		Long:  "",
		RunE:  pullStack,
	}
)

func pullStack(c *cobra.Command, args []string) error {
	log.Info().Msg("TBD")
	return nil
}
