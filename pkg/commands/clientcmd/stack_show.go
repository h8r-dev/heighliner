package clientcmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	stackShowCmd = &cobra.Command{
		Use:   "show [NAME]",
		Short: "Show the description of a stack",
		Long:  "",
		RunE:  showStacks,
	}
)

func showStacks(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please specify the name of the stack")
	}
	for _, stackName := range args {
		showStack(stackName)
	}
	return nil
}

func showStack(name string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Input\tType\tRequired\tDescription")
	for _, s := range defaultStacks {
		if s.Name != name {
			continue
		}
		for _, in := range s.Inputs {
			line := fmt.Sprintf("%s\t%s\t%t\t%s\t", in.Name, in.Type, in.Required, in.Description)
			fmt.Fprintln(w, line)
		}
	}
}
