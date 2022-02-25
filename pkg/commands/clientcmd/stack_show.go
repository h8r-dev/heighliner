package clientcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/h8r-dev/heighliner/pkg/datastore"
	"github.com/spf13/cobra"
)

var (
	stackShowCmd = &cobra.Command{
		Use:   "show [NAME]",
		Short: "Show the description of a stack",
		Long:  "",
		RunE:  showStack,
	}
)

func showStack(c *cobra.Command, args []string) error {
	ds, err := datastore.Stat()
	if err != nil {
		return err
	}
	s, err := ds.Find()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
	err = s.Load()
	if err != nil {
		return err
	}
	line := fmt.Sprintf("%s\t%s\t%s\t", s.Name, s.Version, s.Description)
	fmt.Fprintln(w, line)
	return nil
}
