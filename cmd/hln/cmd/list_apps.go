package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newListAppsCmd() *cobra.Command {
	listAppsCmd := &cobra.Command{
		Use:   "apps",
		Short: "List all heighliner applications",
		Args:  cobra.NoArgs,
	}

	listAppsCmd.RunE = func(c *cobra.Command, args []string) error {

		st, err := getStateInSpecificBackend()
		if err != nil {
			return err
		}

		apps, err := st.ListApps()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
		defer func() {
			err := w.Flush()
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
		}()
		fmt.Fprintln(w, "NAME")
		for _, name := range apps {
			line := fmt.Sprintf("%s\t", name)
			fmt.Fprintln(w, line)
		}
		return nil
	}

	return listAppsCmd
}
