package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
)

type showOptions struct {
	Stack   string
	Version string
	Dir     string

	genericclioptions.IOStreams
}

func (o *showOptions) Validate(cmd *cobra.Command, args []string) error {
	errs := validation.IsDNS1123Subdomain(args[0])
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}

func (o *showOptions) Complete(stackName string) error {
	o.Stack = stackName
	if strings.Contains(o.Stack, "@") {
		args := strings.Split(o.Stack, "@")
		if len(args) < 2 {
			return errors.New("invalid stack fromat, should be name@version")
		}
		o.Stack = args[0]
		o.Version = args[1]
	}
	return nil
}

func (o *showOptions) Run(stackName string) error {
	stk, err := stack.New(o.Stack, o.Version)
	if err != nil {
		return err
	}
	if err := stk.Update(); err != nil {
		return err
	}
	o.Dir = stk.Path
	meta, err := stack.LoadMeta(o.Dir)
	if err != nil {
		return err
	}
	meta.Show(o.Out)
	schema := schema.New(o.Dir)
	if err := schema.LoadSchema(); err != nil {
		return err
	}
	schema.Show(o.Out)
	return nil
}

func newShowCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &showOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:     "show [stack]",
		Short:   "show metadata of stack",
		Aliases: []string{"describe"},
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.Validate(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(args[0]); err != nil {
				return err
			}
			return o.Run(args[0])
		},
	}
	return cmd
}
