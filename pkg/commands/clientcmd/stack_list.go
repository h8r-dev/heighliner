package clientcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/spf13/cobra"
)

var (
	stackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all stacks",
		Args:  cobra.NoArgs,
		RunE:  listStack,
	}
)

func listStack(c *cobra.Command, args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
	for _, v := range stack.Stacks {
		line := fmt.Sprintf("%s\t%s\t%s\t", v.Name, v.Version, v.Description)
		fmt.Fprintln(w, line)
	}
	return nil
}
