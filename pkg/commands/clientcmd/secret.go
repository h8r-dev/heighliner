package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	secretCmd = &cobra.Command{
		Use:   "secret",
		Short: "Manage secrets",
		Long:  "",
	}
)
