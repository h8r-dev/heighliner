package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/stack"
)

func newListStacksCmd(streams genericclioptions.IOStreams) *cobra.Command {
	listStacksCmd := &cobra.Command{
		Use:   "stacks",
		Short: "List stacks",
		Args:  cobra.NoArgs,
	}

	listStacksCmd.RunE = func(c *cobra.Command, args []string) error {
		ss, err := stack.List()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(streams.Out, 0, 4, 2, ' ', tabwriter.TabIndent)
		defer func() {
			err := w.Flush()
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
		}()
		fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
		for _, s := range ss {
			line := fmt.Sprintf("%s\t%s\t%s", s.Name, s.Version, s.Description)
			fmt.Fprintln(w, line)
		}
		return nil
	}

	return listStacksCmd
}
