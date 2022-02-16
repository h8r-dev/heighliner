package clientcmd

import (
	"github.com/spf13/cobra"
)

var (
	stackhubCmd = &cobra.Command{
		Use:   "stackhub",
		Short: "Manage stackhubs",
		Long:  "",
	}
)

type stackhub struct {
	Repo   string
	Branch string
	Path   string
}
