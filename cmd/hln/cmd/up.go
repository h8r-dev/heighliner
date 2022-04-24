package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
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

const upDesc = `
This command run a stack.

You should use '-s' or '--stack' to specify the stack. Use 'list stacks' subcommand 
to check all available stacks. Alternatively, you can use '--dir' flag 
to specify a local directory as your stack source. If you don't specify both '-s' 
and '--dir' flag, it will use current working directory by default:

    $ hln up -s gin-next

or

    $ hln up --dir /path/to/your/stack

To set values in a stack, use '-s' or '--stack' flag to specify a stack, use 
the '--set' flag and pass configuration from the command line:

    $ hln up -s gin-next --set foo=bar

You can specify the '--set' flag multiple times. The priority will be given to the
last (right-most) set specified. For example, if both 'bar' and 'newbar' values are
set for a key called 'foo', the 'newbar' value would take precedence:

    $ hln up -s gin-next --set foo=bar --set foo=newbar

Simply set '-i' or '--interactive' flag and it will run the stack interactively. You can 
fill your input values according to the prompts:

    $ hln up -s gin-next -i

`

// upOptions controls the behavior of up command.
type upOptions struct {
	Stack string
	Dir   string
	local bool

	Values []string

	Interactive bool
	NoCache     bool

	genericclioptions.IOStreams
}

func (o *upOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.Stack, "stack", "s", "", "Name of your stack")
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVarP(&o.Interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *upOptions) Validate(cmd *cobra.Command, args []string) error {
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

func (o *upOptions) Run() error {
	if o.Dir != "" {
		var err error
		o.Dir, err = homedir.Expand(o.Dir)
		if err != nil {
			return err
		}
		o.local = true
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
		p = project.New(pwd, filepath.Join(state.GetTemp(), filepath.Base("hln")))
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
		envvar[1], err = homedir.Expand(envvar[1])
		if err != nil {
			return err
		}
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
		Name:    "up",
		Dir:     o.Dir,
		Plan:    "./plans",
		NoCache: o.NoCache,
	})
	if err != nil {
		return err
	}

	// Print the output.
	b, err := os.ReadFile(stackOutput)
	if err != nil {
		return err
	}
	fmt.Fprintf(o.IOStreams.Out, "%s\n", b)

	// Keep the output info.
	err = copy.Copy(stackOutput, filepath.Join(pwd, appInfo))
	if err != nil {
		return err
	}

	return nil
}

func newUpCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &upOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up your application",
		Long:  upDesc,
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
