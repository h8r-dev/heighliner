package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/project"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
)

// upOptions controls the behavior of up command.
type downOptions struct {
	Stack string
	Dir   string
	Local bool

	Values []string

	Interactive bool
	NoCache     bool

	genericclioptions.IOStreams
}

func (o *downOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.Stack, "stack", "s", "", "Name of your stack")
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVarP(&o.Interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *downOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Stack != "" && o.Dir != "" {
		return errors.New("please do not specify both stack and dir")
	}
	for _, v := range o.Values {
		if !strings.Contains(v, "=") {
			return errors.New("value format should be '--set key=value'")
		}
	}
	return nil
}

func (o *downOptions) Run() error {
	if o.Dir != "" {
		var err error
		o.Dir, err = homedir.Expand(o.Dir)
		if err != nil {
			return err
		}
		o.Local = true
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Prepare the satck
	var p *project.Project
	switch {
	case o.Stack != "":
		s, err := stack.New(o.Stack)
		if err != nil {
			return err
		}
		err = s.Update()
		if err != nil {
			return err
		}
		p = project.New(
			filepath.Join(state.GetCache(), s.Name),
			filepath.Join(state.GetTemp(), s.Name))

	case o.Dir != "":
		p = project.New(o.Dir, filepath.Join(state.GetTemp(), filepath.Base(o.Dir)))
	default:
		p = project.New(pwd, filepath.Join(state.GetTemp(), "hln"))
	}

	// Initialize the project.
	// Enter the project dir automatically.
	err = p.Init()
	if err != nil {
		return err
	}

	// Set input values.
	for _, val := range o.Values {
		envvar := strings.Split(val, "=")
		err := os.Setenv(envvar[0], envvar[1])
		if err != nil {
			return err
		}
	}
	if o.Interactive {
		sch := schema.New()
		err = sch.AutomaticEnv(o.Interactive)
		if err != nil {
			return err
		}
	}

	// Execute the action.
	cli, err := dagger.NewClient(
		viper.GetString("log-format"),
		viper.GetString("log-level"),
		o.IOStreams,
	)
	if err != nil {
		return err
	}
	err = cli.Do(&dagger.ActionOptions{
		Name:    "down",
		Dir:     o.Dir,
		Plan:    "./plans",
		NoCache: o.NoCache,
	})
	if err != nil {
		return err
	}

	return nil
}

func newDownCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &downOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Take down your application",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(cmd, args); err != nil {
				return err
			}
			return o.Run()
		},
	}
	o.BindFlags(cmd.Flags())

	return cmd
}
