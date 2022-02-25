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

// Currently not used. Comment out to pass ci test
// type stackhub struct {
// 	Repo   string
// 	Branch string
// 	Path   string
// }
