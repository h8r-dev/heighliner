package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
)

// testOptions controls the behavior of test command.
type testOptions struct {
	Dir  string
	Plan string

	Values []string

	NoCache bool

	genericclioptions.IOStreams
}

func (o *testOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.StringVarP(&o.Plan, "plan", "p", "./", "Relative path to your test plan")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *testOptions) Validate(cmd *cobra.Command, args []string) error {
	for _, v := range o.Values {
		if !strings.Contains(v, "=") {
			return errors.New("value format should be '--set key=value'")
		}
	}
	return nil
}

func (o *testOptions) Run() error {
	// Set input values.
	for _, val := range o.Values {
		envvar := strings.Split(val, "=")
		err := os.Setenv(envvar[0], envvar[1])
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
		Name:    "test",
		Dir:     o.Dir,
		Plan:    o.Plan,
		NoCache: o.NoCache,
	})
	if err != nil {
		return err
	}

	return nil
}

func newTestCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := testOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test your stack",
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
