package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/stack"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List input values",
		Args:  cobra.NoArgs,
		RunE:  listStacks,
	}
)

func listStacks(c *cobra.Command, args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer func() {
		err := w.Flush()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()
	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
	for _, v := range stack.Stacks {
		line := fmt.Sprintf("%s\t%s\t%s\t", v.Name, v.Version, v.Description)
		fmt.Fprintln(w, line)
	}
	return nil
}
