package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		Long:  "",
	}
)
