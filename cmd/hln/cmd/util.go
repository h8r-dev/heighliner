package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/rwtodd/Go.Sed/sed"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newUtilCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "util",
		Short:  "utils of hln",
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	cmd.AddCommand(newUtilSedCmd(streams))

	return cmd
}

func newUtilSedCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sed",
		Short: "stream editor for filtering and transforming text",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buf := bytes.NewBuffer([]byte("s?server: https://.*?server: https://kubernetes.default.svc?"))
			engine, err := sed.New(buf)
			if err != nil {
				return err
			}
			rawPath := args[0]
			expandPath, err := homedir.Expand(rawPath)
			if err != nil {
				return err
			}
			b, err := ioutil.ReadFile(expandPath)
			if err != nil {
				return err
			}
			str, err := engine.RunString(string(b))
			if err != nil {
				return err
			}
			fmt.Fprint(streams.Out, str)
			return nil
		},
	}
	return cmd
}
