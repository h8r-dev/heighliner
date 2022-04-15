package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/project"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/dagger"
)

// upOptions controls the behavior of up command.
type upOptions struct {
	Stack string
	Path  string

	Values []string

	Interactive bool
	NoCache     bool
}

func (o *upOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.Stack, "stack", "s", "", "Name of your stack")
	f.StringVar(&o.Path, "dir", "", "Path to your local stack")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVarP(&o.Interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *upOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Stack != "" && o.Path != "" {
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
			path.Join(state.GetCache(), s.Name),
			path.Join(state.GetTemp(), s.Name))

	case o.Path != "":
		sp, err := homedir.Expand(o.Path)
		if err != nil {
			return err
		}
		p = project.New(sp, filepath.Join(state.GetTemp(), filepath.Base(sp)))
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
	newArgs := []string{}
	newArgs = append(newArgs,
		"--log-format", viper.GetString("log-format"),
		"--log-level", viper.GetString("log-level"),
		"do", "up",
		"-p", "./plans")
	if o.NoCache {
		newArgs = append(newArgs, "--no-cache")
	}
	err = util.Exec(
		dagger.GetPath(),
		newArgs...)
	if err != nil {
		return err
	}

	// Print the output.
	b, err := os.ReadFile(stackOutput)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s\n", b)

	// Keep the output info.
	err = copy.Copy(stackOutput, filepath.Join(pwd, appInfo))
	if err != nil {
		return err
	}

	return nil
}

func newUpCmd() *cobra.Command {
	o := &upOptions{}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up your application",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			lg := logger.New()
			err := o.Validate(cmd, args)
			if err != nil {
				lg.Fatal().Err(err).Msg("invalid args")
			}
			err = o.Run()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to run")
			}
		},
	}
	o.BindFlags(cmd.Flags())

	return cmd
}
