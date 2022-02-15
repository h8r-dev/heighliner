package clientcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	stackListCmd = &cobra.Command{
		Use:   "list",
		Short: "List stacks",
		Long:  "",
		RunE:  listStacks,
	}
)

func listStacks(c *cobra.Command, args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
	for _, s := range defaultStacks {
		line := fmt.Sprintf("%s\t%s\t%s\t", s.Name, s.Version, s.Description)
		fmt.Fprintln(w, line)
	}
	return nil
}
